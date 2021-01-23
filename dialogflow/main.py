from flask import Flask, redirect, url_for, request, jsonify
from dflow import getApiAiAnswer

app = Flask(__name__)
app.config["DEBUG"] = True


@app.route('/api/v1/getAnswer', methods=['POST'])
def getAnswer():
    # req data struct :
    # message = ""
    data = request.get_json()
    try:
        if data["message"] == "":
            return jsonify({"err":"message field is empty"})
    except:
        return jsonify({"err":"no message field"})

    message = getApiAiAnswer(data["message"])
    if message == "":
        return jsonify({"err":"error making req"})
    
    return jsonify({"answer":message})

if __name__ == "__main__":
    app.run()
