## Trak - Learning Tool Locally For Learners 

## Common Flow 
1. User May write `kavro init lang/go` | `kavro init tool/jenkins`
2. My Engine First Search into my Registry with Available Tools 
    - Steps:
        i. Finding `lang/go` | `tool/jenkins`
        ii. Found, Fetching `lang/go` | `tool/jenkins` ....
        iii. Find where to put the directories and subdiretories 
            - May give path as well like `kavro init lang/go -path=C:/users/username/desktop/`
        iv. Fetch Actual Data For `lang/go` | `tool/jenkins` (can be in json, yaml, toml) i think json 
        v.  Efficiently Each structured data in sequence , create folders, files sequenctially one by one , also kavro.json 
        vi. thats it , for now 
3. User Will Go Through the topics now 

### Need to fix 
1. where from the finding whether that lang exist or not
2. we can integrate (finding + fetching) in one take
3. where from the data for the specific `lang/go` | `tool/jenkins` will come frome, because it will be having big json data u know, (githhub raw one ?, my own , if i take my own i dont have money so i will take help of github itself) 
i will kind of add `kavro_availables.json` it will have what i have in my system, and `kavro_learn_go.json`, ..... 
4. In Each json what data will be there ?