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

type VersionConfig struct {
    Name        string       `json:"name"`
    DisplayName string       `json:"display_name"`
    Files       []FileConfig `json:"files"`
    Profile     struct {
        URL string `json:"url"`
    } `json:"profile"`
}

type Config struct {
    Versions []VersionConfig `json:"versions"`
}

func main() {
    fmt.Println("==========================================")
    fmt.Println("     YAYO Minecraft 모드팩 설치기")
    fmt.Println("==========================================")
    fmt.Println()

    // 설정 파일 다운로드
    fmt.Println("=== 설정 파일 다운로드 중... ===")
    configURL := "https://files.yayokorea.net/share/minecraft/config.json"
    resp, err := http.Get(configURL)
    if err != nil {
        panic(fmt.Sprintf("설정 파일 다운로드 실패: %v", err))
    }
    defer resp.Body.Close()

    var config Config
    if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
        panic(fmt.Sprintf("설정 파일 파싱 실패: %v", err))
    }
    fmt.Println("   ✓ 설정 파일 다운로드 완료")
    fmt.Println()

    // 버전 선택
    fmt.Println("=== 설치할 버전을 선택하세요 ===")
    for i, version := range config.Versions {
        fmt.Printf("%d. %s\n", i+1, version.DisplayName)
    }
    fmt.Print("\n선택 (1-", len(config.Versions), "): ")

    var choice int
    for {
        _, err := fmt.Scanf("%d", &choice)
        if err != nil {
            fmt.Print("잘못된 입력입니다. 다시 선택해주세요: ")
            continue
        }
        if choice < 1 || choice > len(config.Versions) {
            fmt.Print("잘못된 번호입니다. 다시 선택해주세요: ")
            continue
        }
        break
    }

    selectedVersion := config.Versions[choice-1]
    fmt.Printf("\n'%s' 설치를 시작합니다.\n\n", selectedVersion.DisplayName)

    // 환경 변수 처리
    appData := os.Getenv("APPDATA")
    if appData == "" {
        panic("APPDATA 환경 변수를 찾을 수 없습니다.")
    }

    // launcher_profiles.json 수정
    fmt.Println("=== Minecraft 프로파일 설정 ===")
    fmt.Println("   launcher_profile.json 수정 중...")

    launcherProfiles := filepath.Join(appData, ".minecraft", "launcher_profiles.json")
    updateProfile(launcherProfiles, selectedVersion.Profile.URL, selectedVersion.Name)

    fmt.Println("   ✓ 프로파일 설정 완료")
    fmt.Println()

    // 파일 다운로드 및 설치
    tempDir := filepath.Join(os.Getenv("TEMP"), "YAYO")
    os.MkdirAll(tempDir, os.ModePerm)

    fmt.Println("=== 파일 다운로드 및 설치 ===")
    for _, file := range selectedVersion.Files {
        // 환경 변수 치환
        dest := strings.Replace(file.Destination, "%APPDATA%", appData, -1)

        // 파일명 추출
        fileName := filepath.Base(file.URL)
        fmt.Printf("   %s 다운로드 중...\n", fileName)

        // 다운로드
        tempFile := filepath.Join(tempDir, fileName)
        downloadFile(file.URL, tempFile)
        fmt.Println()

        // 압축 해제 전 ZIP 파일 포맷 확인
        fmt.Printf("   %s 압축 해제 중...\n", fileName)
        if !isZipFile(tempFile) {
            panic("다운로드된 파일이 ZIP 포맷이 아닙니다. 서버 파일 또는 URL을 확인하세요.")
        }
		info, err := os.Stat(tempFile)
		if err != nil {
			panic(fmt.Sprintf("임시 파일 정보 확인 실패: %v", err))
		}
		fmt.Printf("압축 해제 대상 파일: %s, 크기: %d bytes\n", tempFile, info.Size())
        unzip(tempFile, dest)
        fmt.Println()

        // 임시 파일 삭제
        os.Remove(tempFile)
        fmt.Printf("   ✓ %s 삭제 완료\n", fileName)
        fmt.Println()
    }
    fmt.Println()

    fmt.Println("==========================================")
    fmt.Println("              설치 완료!")
    fmt.Println("    모든 작업이 성공적으로 완료되었습니다")
    fmt.Println("==========================================")
    fmt.Println()

    fmt.Println("종료하려면 Enter 키를 누르세요...")
    var input string
    fmt.Scanln(&input)
    fmt.Scanln(&input)
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

    // 파일 크기 출력
    fileSize := resp.ContentLength
    if fileSize >= 1048576 {
        fmt.Printf("   ▶ 파일 크기: %.2f MB\n", float64(fileSize)/1048576)
    } else {
        fmt.Printf("   ▶ 파일 크기: %.2f KB\n", float64(fileSize)/1024)
    }

    // 전체 파일 복사
    written, err := io.Copy(out, resp.Body)
    if err != nil {
        panic(fmt.Sprintf("파일 저장 실패: %v", err))
    }

    // 다운로드 완료 출력
    if written >= 1048576 {
        fmt.Printf("   ✓ 다운로드 완료: %.2f MB\n", float64(written)/1048576)
    } else {
        fmt.Printf("   ✓ 다운로드 완료: %.2f KB\n", float64(written)/1024)
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

// ZIP 파일 포맷 확인 함수
func isZipFile(path string) bool {
    file, err := os.Open(path)
    if err != nil {
        return false
    }
    defer file.Close()
    header := make([]byte, 2)
    _, err = file.Read(header)
    if err != nil {
        return false
    }
    return header[0] == 'P' && header[1] == 'K'
}

func updateProfile(launcherProfiles string, profileURL string, version string) {
    var jsonData map[string]interface{}

    // launcher_profiles.json 파일이 없으면 기본 구조로 생성
    if _, err := os.Stat(launcherProfiles); os.IsNotExist(err) {
        jsonData = map[string]interface{}{
            "profiles":       map[string]interface{}{},
            "settings":       map[string]interface{}{},
            "version_groups": []interface{}{},
        }
        // 디렉토리가 없다면 생성
        if err := os.MkdirAll(filepath.Dir(launcherProfiles), os.ModePerm); err != nil {
            panic(fmt.Sprintf("디렉토리 생성 실패: %v", err))
        }
    } else {
        // 파일이 있으면 읽기
        file, err := os.Open(launcherProfiles)
        if err != nil {
            panic(fmt.Sprintf("launcher_profiles.json 열기 실패: %v", err))
        }
        defer file.Close()

        decoder := json.NewDecoder(file)
        if err := decoder.Decode(&jsonData); err != nil {
            panic(fmt.Sprintf("launcher_profiles.json 파싱 실패: %v", err))
        }
    }

    // 프로파일 다운로드 및 추가
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

    // 파일 저장
    launcherFile, err := os.Create(launcherProfiles)
    if err != nil {
        panic(fmt.Sprintf("launcher_profiles.json 쓰기 실패: %v", err))
    }
    defer launcherFile.Close()

    encoder := json.NewEncoder(launcherFile)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(jsonData); err != nil {
        panic(fmt.Sprintf("launcher_profiles.json 업데이트 실패: %v", err))
    }
}