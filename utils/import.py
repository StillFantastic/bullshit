# Help you import data

import json

with open('DataToAdd.txt', 'r', encoding="utf-8") as f:
    Add_Data = []
    for temp_data in f.readlines():
        Add_Data.append(temp_data.strip())

with open('../generator/data.json', 'r', encoding='utf-8') as fi:
    ori_data = json.loads(fi.read())
    # famous, bullshit, before, after
    ori_data['famous'].extend(Add_Data)

json_str = json.dumps(ori_data, indent=4, ensure_ascii=False)

with open('../generator/data.json', 'w', encoding='utf-8') as json_file:
    json_file.write(json_str)
