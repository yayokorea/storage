package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// 환경 변수에서 %APPDATA% 경로 가져오기
	appData := os.Getenv("APPDATA")
	if appData == "" {
		panic("APPDATA 환경 변수를 찾을 수 없습니다.")
	}

	// 경로 설정
	tempDir := "C:\\Windows\\Temp\\YAYO"
	minecraftDir := filepath.Join(appData, ".minecraft")
	launcherProfiles := filepath.Join(minecraftDir, "launcher_profiles.json")
	minecraftBaseDir := "C:\\Users\\Public\\Minecraft"
	github := "https://raw.githubusercontent.com/yayokorea/storage/main/"
	storage := "1.21_Fabric/"
	file1 := "1.21_Fabric_Survival.zip"
	file2 := "YAYO_1.21.zip"
	version := "YAYO_1.21"

	// 임시 디렉터리 생성
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
			panic(fmt.Sprintf("임시 디렉터리 생성 실패: %v", err))
		}
	}

	// launcher_profiles.json 수정
	fmt.Println("launcher_profile.json 수정 중")
	if _, err := os.Stat(launcherProfiles); os.IsNotExist(err) {
		panic(fmt.Sprintf("launcher_profiles.json 파일을 찾을 수 없습니다: %s", launcherProfiles))
	}

	var jsonData map[string]interface{}
	file, err := os.Open(launcherProfiles)
	if err != nil {
		panic(fmt.Sprintf("launcher_profiles.json 열기 실패: %v", err))
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonData); err != nil {
		panic(fmt.Sprintf("launcher_profiles.json 파싱 실패: %v", err))
	}

	profileURL := fmt.Sprintf("%s%sprofile.json", github, storage)
	resp, err := http.Get(profileURL)
	if err != nil {
		panic(fmt.Sprintf("프로파일 다운로드 실패: %v", err))
	}
	defer resp.Body.Close()

	var profileData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&profileData); err != nil {
		panic(fmt.Sprintf("프로파일 JSON 파싱 실패: %v", err))
	}

	profiles := jsonData["profiles"].(map[string]interface{})
	profiles[version] = profileData

	launcherFile, err := os.Create(launcherProfiles)
	if err != nil {
		panic(fmt.Sprintf("launcher_profiles.json 쓰기 실패: %v", err))
	}
	defer launcherFile.Close()

	encoder := json.NewEncoder(launcherFile)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(jsonData); err != nil {
		panic(fmt.Sprintf("launcher_profiles.json 업데이트 실패: %v", err))
	}

	// 파일 다운로드
	fmt.Printf("%s 다운로드 중\n", file1)
	downloadFile(fmt.Sprintf("%s%s%s", github, storage, file1), filepath.Join(tempDir, file1))
	fmt.Printf("%s 다운로드 중\n\n", file2)
	downloadFile(fmt.Sprintf("%s%s%s", github, storage, file2), filepath.Join(tempDir, file2))

	// 압축 해제
	fmt.Printf("%s 압축 해제 중\n", file1)
	unzip(filepath.Join(tempDir, file1), minecraftBaseDir)
	fmt.Printf("%s 압축 해제 중\n\n", file2)
	unzip(filepath.Join(tempDir, file2), filepath.Join(minecraftDir, "versions"))

	// 다운로드한 파일 삭제
	fmt.Printf("%s 삭제 중\n", file1)
	os.Remove(filepath.Join(tempDir, file1))
	fmt.Printf("%s 삭제 중\n\n", file2)
	os.Remove(filepath.Join(tempDir, file2))

	// 설치 완료 메시지
	fmt.Println("설치 완료!")

	// 종료 전에 Enter 키를 기다리기
	fmt.Println("계속하려면 Enter 키를 누르세요...")
	fmt.Scanln() // 사용자 입력을 기다립니다.
}

func downloadFile(url string, dest string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("파일 다운로드 실패: %v", err))
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		panic(fmt.Sprintf("파일 생성 실패: %v", err))
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		panic(fmt.Sprintf("파일 쓰기 실패: %v", err))
	}
}

func unzip(src string, dest string) {
	// 압축 파일 열기
	r, err := zip.OpenReader(src)
	if err != nil {
		panic(fmt.Sprintf("압축 파일 열기 실패: %v", err))
	}
	defer r.Close()

	// 압축 파일 안의 폴더 이름을 추출
	var folderName string
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			folderName = f.Name
			break
		}
	}

	// 해당 폴더가 이미 존재하면 삭제
	if folderName != "" {
		destFolder := filepath.Join(dest, folderName)
		if _, err := os.Stat(destFolder); !os.IsNotExist(err) {
			fmt.Printf("%s 폴더 삭제 중\n", destFolder)
			if err := os.RemoveAll(destFolder); err != nil {
				panic(fmt.Sprintf("폴더 삭제 실패: %v", err))
			}
		}
	}

	// 새로운 디렉터리 생성
	if err := os.MkdirAll(filepath.Join(dest, folderName), os.ModePerm); err != nil {
		panic(fmt.Sprintf("디렉터리 생성 실패: %v", err))
	}

	// 압축 해제
	for _, f := range r.File {
		fPath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			panic(fmt.Sprintf("디렉터리 생성 실패: %v", err))
		}

		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(fmt.Sprintf("파일 생성 실패: %v", err))
		}

		rc, err := f.Open()
		if err != nil {
			panic(fmt.Sprintf("압축 파일 읽기 실패: %v", err))
		}

		if _, err := io.Copy(outFile, rc); err != nil {
			panic(fmt.Sprintf("파일 쓰기 실패: %v", err))
		}

		outFile.Close()
		rc.Close()
	}
}