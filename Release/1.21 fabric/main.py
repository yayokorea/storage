#pyinstaller --uac-admin  --icon="C:\Users\user\Documents\NSIS\icon.ico" "C:\Users\user\Desktop\Github\MinecraftInstaller\storage\Release\1.21 fabric\main.py"
import os
from requests import get
import zipfile
import json
import requests
import shutil

temp = "C:/Windows/Temp/YAYO/"
minecraft = f"{os.getenv('APPDATA')}/.minecraft/"

github = 'https://raw.githubusercontent.com/yayokorea/storage/main/'
storage = '1.21_Fabric/'

file1 = '1.21_Fabric_Survival.zip'
file2 = "YAYO_1.21.zip"
version = "YAYO_1.21"

if os.path.isdir(temp):
    pass
else:
    os.mkdir(temp)

if __name__ == '__main__':

    if os.path.isdir('C:/Users/Public/Minecraft/1.21_Fabric_Survival/'):
        shutil.rmtree('C:/Users/Public/Minecraft/1.21_Fabric_Survival/')
    else:
        pass

    print('launcher_profile.json 수정중\n')
    with open(f'{minecraft}launcher_profiles.json', 'r', encoding='utf-8') as f:
        json_data = json.load(f)
        url = f'{github}{storage}profile.json'
        data = requests.get(url).json()
        json_data['profiles'][version] = data

    with open(f'{minecraft}launcher_profiles.json', 'w', encoding='utf-8') as f:
        json.dump(json_data, f, indent="\t")

    def download(url, file_name):
        with open(file_name, "wb") as file:  
            response = get(url)            
            file.write(response.content)    

    url_1 = f"{github}{storage}{file1}"
    url_2 = f"{github}{storage}{file2}"
    
    print(f"{file1} 다운로드중")
    download(url_1,f"{temp}{file1}")
    print(f"{file2} 다운로드중\n")
    download(url_2,f"{temp}{file2}")

    print(f"{file1} 압축해제중")
    zipfile.ZipFile(f'{temp}{file1}').extractall(path='C:/Users/Public/Minecraft/')
    print(f"{file2} 압축해제중\n")
    zipfile.ZipFile(f'{temp}{file2}').extractall(path=f'{minecraft}versions/')
    
    print(f"{file1} 삭제중")
    os.remove(f'{temp}{file1}')
    print(f"{file2} 삭제중\n")
    os.remove(f'{temp}{file2}')

    print('설치 완료!\n')
    os.system('pause')

    