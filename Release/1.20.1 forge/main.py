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
    
    url = 'https://raw.githubusercontent.com/yayokorea/storage/main/1.20.1_Forge_Twilight_Forest/profile.json'
    data = requests.get(url).json()
    
    json_data['profiles']["1.20.1_YAYO"] = data

with open(f'{minecraft}launcher_profiles.json', 'w', encoding='utf-8') as f:
    json.dump(json_data, f, indent="\t")

def download(url, file_name):
    with open(file_name, "wb") as file:  
        response = get(url)            
        file.write(response.content)     

if __name__ == '__main__':
    url_1 = "https://raw.githubusercontent.com/yayokorea/storage/main/1.20.1_Forge_Twilight_Forest/1.20.1_YAYO_Forge.zip"
    url_2 = "https://raw.githubusercontent.com/yayokorea/storage/main/1.20.1_Forge_Twilight_Forest/YAYO_1.20.1.zip"
    url_3 = "https://raw.githubusercontent.com/yayokorea/storage/main/1.20.1_Forge_Twilight_Forest/minecraftforge.zip"
    url_4 = "https://raw.githubusercontent.com/yayokorea/storage/main/1.20.1_Forge_Twilight_Forest/minecraft.zip"
    
    
    print("1.20.1_YAYO_Forge.zip 설치중")
    download(url_1,f"{temp}1.20.1_YAYO_Forge.zip")
    print("YAYO_1.20.1.zip 설치중")
    download(url_2,f"{temp}YAYO_1.20.1.zip")
    print("minecraftforge.zip 설치중")
    download(url_3,f"{temp}minecraftforge.zip")
    print("minecraft.zip 설치중\n")
    download(url_4,f"{temp}minecraft.zip")
    

    print("1.20.1_YAYO_Forge.zip 압축해제중")
    zipfile.ZipFile(f'{temp}1.20.1_YAYO_Forge.zip').extractall(path='C:/Users/Public/Minecraft/')
    print("YAYO_1.20.1.zip 압축해제중")
    zipfile.ZipFile(f'{temp}YAYO_1.20.1.zip').extractall(path=f'{minecraft}versions/')
    print("minecraftforge.zip 압축해제중")
    zipfile.ZipFile(f'{temp}minecraftforge.zip').extractall(path=f'{minecraft}libraries/net/')
    print("minecraft.zip 압축해제중\n")
    zipfile.ZipFile(f'{temp}minecraft.zip').extractall(path=f'{minecraft}libraries/net/')
    
    

    print("1.20.1_YAYO_Forge.zip 삭제중")
    os.remove(f'{temp}1.20.1_YAYO_Forge.zip')
    print("YAYO_1.20.1.zip 삭제중")
    os.remove(f'{temp}YAYO_1.20.1.zip')
    print("minecraftforge.zip 삭제중")
    os.remove(f'{temp}minecraftforge.zip')
    print("minecraft.zip 삭제중\n")
    os.remove(f'{temp}minecraft.zip')

    print('설치 완료!\n')
    os.system('pause')