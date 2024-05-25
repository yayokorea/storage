#pyinstaller --uac-admin  --icon="C:\Users\user\Documents\NSIS\icon.ico" "C:\Users\user\Desktop\cxfreeze\main.py"
import os
from requests import get
import zipfile
import json
import requests

temp = "C:/Windows/Temp/YAYO/"
minecraft = f"{os.getenv('APPDATA')}/.minecraft/"

if os.path.isdir(temp):
    pass
else:
    os.mkdir(temp)


with open(f'{minecraft}launcher_profiles.json', 'r', encoding='utf-8') as f:
    json_data = json.load(f)
    
    url = 'https://raw.githubusercontent.com/yayokorea/storage/main/1.20.6_Fabric/profile.json'
    data = requests.get(url).json()
    
    json_data['profiles']["1.20.1_YAYO"] = data

with open(f'{minecraft}launcher_profiles.json', 'w', encoding='utf-8') as f:
    json.dump(json_data, f, indent="\t")

def download(url, file_name):
    with open(file_name, "wb") as file:  
        response = get(url)            
        file.write(response.content)     

if __name__ == '__main__':
    url_1 = "https://raw.githubusercontent.com/yayokorea/storage/main/1.20.6_Fabric/1.20.6%20Fabric%20Survival.zip"
    url_2 = "https://raw.githubusercontent.com/yayokorea/storage/main/1.20.6_Fabric/YAYO_1.20.6.zip"
    
    
    print("1.20.6 Fabric Survival.zip 설치중")
    download(url_1,f"{temp}1.20.6 Fabric Survival.zip")
    print("YAYO_1.20.6.zip 설치중")
    download(url_2,f"{temp}YAYO_1.20.6.zip")

    print("1.20.6 Fabric Survival.zip 압축해제중")
    zipfile.ZipFile(f'{temp}1.20.6 Fabric Survival.zip').extractall(path='C:/Users/Public/Minecraft/')
    print("YAYO_1.20.6.zip 압축해제중")
    zipfile.ZipFile(f'{temp}YAYO_1.20.6.zip').extractall(path=f'{minecraft}versions/')

    print("1.20.6 Fabric Survival.zip 삭제중")
    os.remove(f'{temp}1.20.6 Fabric Survival.zip')
    print("YAYO_1.20.6.zip 삭제중")
    os.remove(f'{temp}YAYO_1.20.6.zip')

    print('설치 완료!\n')
    os.system('pause')