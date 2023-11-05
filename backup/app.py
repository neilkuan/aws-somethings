import boto3
import json
from datetime import datetime

date = str(datetime.today())
def main():
  policy_list = []
  iam = boto3.client('iam')
  data = iam.list_policies(Scope='AWS', MaxItems=1000)
  for policy in data['Policies']:
    policy_list.append(policy.get('PolicyName'))
  
  before_policy_list = []
  with open('policies.json') as file: 
    before_policy_list = json.load(file)

  if before_policy_list != policy_list:
    print(f"[{date}] Start to Update policies.json...")
    with open('policies.json', 'w') as file:
      file.write(json.dumps(policy_list, indent=4, sort_keys=True))
      
    print(f"[{date}] Updated policies.json...")
    with open('update.txt', 'w') as file:
      file.writelines('Y')
    
  else:
    print(f"[{date}] Not Update policies.json...")
    with open('update.txt', 'w') as file:
      file.writelines('N')
if __name__ == "__main__":
    main()