package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileConfig struct {
	URL         string `json:"url"`
	Destination string `json:"destination"`
}

type Config struct {
	Name    string       `json:"name"`
	Files   []FileConfig `json:"files"`
	Profile struct {
		URL string `json:"url"`
	} `json:"profile"`
}

func main() {
	fmt.Println("==========================================")
	fmt.Println("     YAYO Minecraft 모드팩 설치기")
	fmt.Println("==========================================")
	fmt.Println()

	// 설정 파일 다운로드
	fmt.Println("=== 설정 파일 다운로드 ===")
	configURL := "http://58.231.82.19:8000/config.json"
	resp, err := http.Get(configURL)
	if err != nil {
		panic(fmt.Sprintf("설정 파일 다운로드 실패: %v", err))
	}
	defer resp.Body.Close()

	var config Config
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		panic(fmt.Sprintf("설정 파일 파싱 실패: %v", err))
	}

	// 환경 변수 처리
	appData := os.Getenv("APPDATA")
	if appData == "" {
		panic("APPDATA 환경 변수를 찾을 수 없습니다.")
	}

	// launcher_profiles.json 수정
	fmt.Println("=== Minecraft 프로파일 설정 ===")
	fmt.Println("   launcher_profile.json 수정 중...")

	launcherProfiles := filepath.Join(appData, ".minecraft", "launcher_profiles.json")
	updateProfile(launcherProfiles, config.Profile.URL, config.Name)

	fmt.Println("   ✓ 프로파일 설정 완료")
	fmt.Println()

	// 파일 다운로드 및 설치
	tempDir := filepath.Join(os.Getenv("TEMP"), "YAYO")
	os.MkdirAll(tempDir, os.ModePerm)

	fmt.Println("=== 파일 다운로드 및 설치 ===")
	for _, file := range config.Files {
		// 환경 변수 치환
		dest := strings.Replace(file.Destination, "%APPDATA%", appData, -1)

		// 파일명 추출
		fileName := filepath.Base(file.URL)
		fmt.Printf("   %s 다운로드 중...\n", fileName)

		// 다운로드
		tempFile := filepath.Join(tempDir, fileName)
		downloadFile(file.URL, tempFile)
		fmt.Println()

		// 압축 해제
		fmt.Printf("   %s 압축 해제 중...\n", fileName)
		unzip(tempFile, dest)
		fmt.Println()

		// 임시 파일 삭제
		os.Remove(tempFile)
		fmt.Printf("   • %s 삭제 완료\n", fileName)
	}
	fmt.Println()

	fmt.Println("==========================================")
	fmt.Println("              설치 완료!")
	fmt.Println("    모든 작업이 성공적으로 완료되었습니다")
	fmt.Println("==========================================")
	fmt.Println()

	fmt.Println("계속하려면 Enter 키를 누르세요...")
	fmt.Scanln()
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

	fileSize := resp.ContentLength
	if fileSize >= 1048576 { // 1MB 이상
		fmt.Printf("   ▶ 파일 크기: %.2f MB\n", float64(fileSize)/1048576)
	} else {
		fmt.Printf("   ▶ 파일 크기: %.2f KB\n", float64(fileSize)/1024)
	}

	buffer := make([]byte, 1024)
	downloaded := int64(0)

	for {
		n, err := resp.Body.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(fmt.Sprintf("파일 읽기 실패: %v", err))
		}

		if _, err := out.Write(buffer[:n]); err != nil {
			panic(fmt.Sprintf("파일 쓰기 실패: %v", err))
		}

		downloaded += int64(n)
		progress := float64(downloaded) / float64(fileSize) * 100
		fmt.Printf("\r   진행률: %.0f%%", progress)
	}

	if downloaded >= 1048576 { // 1MB 이상
		fmt.Printf("\n   ✓ 다운로드 완료: %.2f MB\n", float64(downloaded)/1048576)
	} else {
		fmt.Printf("\n   ✓ 다운로드 완료: %.2f KB\n", float64(downloaded)/1024)
	}
}

func unzip(src string, dest string) {
	r, err := zip.OpenReader(src)
	if err != nil {
		panic(fmt.Sprintf("압축 파일 열기 실패: %v", err))
	}
	defer r.Close()

	var folderName string
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			folderName = f.Name
			break
		}
	}

	if folderName != "" {
		destFolder := filepath.Join(dest, folderName)
		if _, err := os.Stat(destFolder); !os.IsNotExist(err) {
			fmt.Println("   ▶ 기존 폴더 삭제 중...")
			fmt.Printf("      %s\n", destFolder)
			if err := os.RemoveAll(destFolder); err != nil {
				panic(fmt.Sprintf("폴더 삭제 실패: %v", err))
			}
			fmt.Println("   ✓ 폴더 삭제 완료")
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

func updateProfile(launcherProfiles string, profileURL string, version string) {
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
}
