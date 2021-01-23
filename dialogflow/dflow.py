import const
import apiai, json
import random

def getApiAiAnswer(message):
    request = apiai.ApiAI(const.APIAITOKEN).text_request()
    request.lang = "ru"
    request.session_id = "zhanna"
    request.query = message
    responseJson = json.loads(request.getresponse().read().decode("utf-8"))
    answer = responseJson["result"]["fulfillment"]["speech"]
    return answer
 