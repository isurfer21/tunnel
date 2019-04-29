import requests
import json

APP_URL = "http://127.0.0.1:9999"


class Utility:

    @staticmethod
    def json_validator(data):
        try:
            json.loads(data)
            return True
        except ValueError as error:
            # print("invalid json: %s" % error)
            return False


class TestCase:
    count = 0

    def session(self):
        self.count += 1
        print(self.count)
        url = APP_URL + '/session'
        print(' url\t›', url)
        headers = {
            'cache-control': "no-cache"
        }
        raw_response = requests.request(
            "POST", url, headers=headers
        )
        # print(' raw\t›', raw_response.text)
        json_response = json.loads(raw_response.text)
        print(' json\t›', json_response)
        if Utility.json_validator(json_response['response']):
            result = json.loads(json_response['response'])
            print(' result\t›', result)
            return result['token']
        return None

    def terminal(self, token):
        if token != None:
            self.count += 1
            print(self.count)
            url = APP_URL + '/terminal'
            print(' url\t›', url)
            payload = {
                "token": token,
                "cmd": "ls"
            }
            headers = {
                'content-type': "application/x-www-form-urlencoded",
                'cache-control': "no-cache"
            }
            raw_response = requests.request(
                "POST", url, data=payload, headers=headers
            )
            # print(' raw\t›', raw_response.text)
            json_response = json.loads(raw_response.text)
            print(' json\t›', json_response)
            if Utility.json_validator(json_response['response']):
                result = json.loads(json_response['response'])
                print(' result\t›', result)
            else:
                print(' error\t›', 'invalid_json_response')
        else:
            print(' error\t›', 'missing_token')


print('TestCase')
testCase = TestCase()
token = testCase.session()
testCase.terminal(token)
print('\nDone!')
