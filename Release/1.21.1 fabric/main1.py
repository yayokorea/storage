#pyinstaller --uac-admin  --icon="C:\Users\user\Documents\NSIS\icon.ico" "C:\Users\user\Desktop\Github\MinecraftInstaller\storage\Release\1.21.1 fabric\main.py"
import os
from requests import get
import zipfile
import json
import requests
import shutil

temp = "C:/Windows/Temp/YAYO/"
minecraft = f"{os.getenv('APPDATA')}/.minecraft/"

github = 'https://raw.githubusercontent.com/yayokorea/storage/main/'
storage = '1.21.1_Fabric/'

file1 = '1.21.1_Creative.zip'
file2 = "fabric-loader-0.16.8-1.21.1.zip"
version = "YAYO_1.21.1"

if os.path.isdir(temp):
    pass
else:
    os.mkdir(temp)

def download(url, file_name):
    try:
        response = get(url, timeout=30)  # 타임아웃 추가
        response.raise_for_status()  # HTTP 오류 체크
        with open(file_name, "wb") as file:
            file.write(response.content)
        return True
    except requests.RequestException as e:
        print(f'다운로드 중 오류 발생: {e}')
        return False

if __name__ == '__main__':
    try:
        # 기존 디렉토리 삭제
        if os.path.isdir('C:/Users/Public/Minecraft/1.21.1_Creative/'):
            try:
                shutil.rmtree('C:/Users/Public/Minecraft/1.21.1_Creative/')
            except PermissionError:
                print('디렉토리 삭제 권한이 없습니다.')
                exit(1)
            except Exception as e:
                print(f'디렉토리 삭제 중 오류 발생: {e}')
                exit(1)

        print('launcher_profile.json 수정중\n')
        try:
            with open(f'{minecraft}launcher_profiles.json', 'r', encoding='utf-8') as f:
                json_data = json.load(f)
                url = f'{github}{storage}profile.json'
                response = requests.get(url)
                if response.status_code == 200:
                    data = response.json()
                    json_data['profiles'][version] = data
                else:
                    print('프로필 데이터를 가져오는데 실패했습니다.')
                    exit(1)
        except FileNotFoundError:
            print('launcher_profiles.json 파일을 찾을 수 없습니다.')
            exit(1)
        except json.JSONDecodeError:
            print('프로필 데이터 형식이 잘못되었습니다.')
            exit(1)
        except requests.RequestException:
            print('서버에 연결할 수 없습니다.')
            exit(1)

        try:
            with open(f'{minecraft}launcher_profiles.json', 'w', encoding='utf-8') as f:
                json.dump(json_data, f, indent="\t")
            print('프로필 저장 완료')
        except IOError:
            print('프로필 저장에 실패했습니다.')
            exit(1)

        # 파일 다운로드
        url_1 = f"{github}{storage}{file1}"
        url_2 = f"{github}{storage}{file2}"
    
        print(f"{file1} 다운로드중")
        if not download(url_1, f"{temp}{file1}"):
            print(f"{file1} 다운로드 실패")
            exit(1)

        print(f"{file2} 다운로드중\n")
        if not download(url_2, f"{temp}{file2}"):
            print(f"{file2} 다운로드 실패")
            exit(1)

        # ZIP 파일 압축 해제
        print(f"{file1} 압축해제중")
        try:
            with zipfile.ZipFile(f'{temp}{file1}') as zf:
                zf.extractall(path='C:/Users/Public/Minecraft/')
        except zipfile.BadZipFile:
            print(f"{file1} 파일이 손상되었습니다.")
            exit(1)
        except PermissionError:
            print("압축 해제 권한이 없습니다.")
            exit(1)

        print(f"{file2} 압축해제중\n")
        try:
            with zipfile.ZipFile(f'{temp}{file2}') as zf:
                zf.extractall(path=f'{minecraft}versions/')
        except zipfile.BadZipFile:
            print(f"{file2} 파일이 손상되었습니다.")
            exit(1)
        except PermissionError:
            print("압축 해제 권한이 없습니다.")
            exit(1)

        # 임시 파일 삭제
        print(f"{file1} 삭제중")
        try:
            os.remove(f'{temp}{file1}')
        except OSError as e:
            print(f"파일 삭제 중 오류 발생: {e}")

        print(f"{file2} 삭제중\n")
        try:
            os.remove(f'{temp}{file2}')
        except OSError as e:
            print(f"파일 삭제 중 오류 발생: {e}")

        print('설치 완료!\n')
        os.system('pause')

    except Exception as e:
        print(f'예상치 못한 오류가 발생했습니다: {e}')
        exit(1)

    