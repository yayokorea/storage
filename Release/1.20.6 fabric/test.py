import json
import os

minecraft = f"{os.getenv('APPDATA')}/.minecraft/"

with open(f'{minecraft}launcher_profiles.json', 'r', encoding='utf-8') as f:
    json_data = json.load(f)

    a = {
        "gameDir" : "C:\\Users\\Public\\Minecraft\\1.20.1_YAYO_Forge",
        "javaArgs" : "-Xmx8G -XX:+UnlockExperimentalVMOptions -XX:+UseG1GC -XX:G1NewSizePercent=20 -XX:G1ReservePercent=20 -XX:MaxGCPauseMillis=50 -XX:G1HeapRegionSize=32M",
        "lastUsed" : "9999-12-31T00:00:00.001Z",
        "lastVersionId" : "YAYO_1.20.1",
        "name" : "YAYO_1.20.1",
        "type" : "custom"
        }
    
    json_data['profiles']["1.20.1_YAYO"] = a

    print(json_data)

    with open(f'{minecraft}launcher_profiles.json', 'w', encoding='utf-8') as f:
        json.dump(json_data, f, indent="\t")