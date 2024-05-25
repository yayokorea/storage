import json
import os
import requests

minecraft = f"{os.getenv('APPDATA')}/.minecraft/"

with open(f'{minecraft}launcher_profiles.json', 'r', encoding='utf-8') as f:
    json_data = json.load(f)

    url = 'https://raw.githubusercontent.com/yayokorea/storage/main/1.20.6_Fabric/a.json'
    data = requests.get(url).json()


    print(data)
    
    json_data['profiles']["1.20.1_YAYO"] = data


    with open(f'{minecraft}launcher_profiles.json', 'w', encoding='utf-8') as f:
        json.dump(json_data, f, indent="\t")