"""Server routes"""

from bottle import route, run, template


@route("/")
def index():
    """Route homepage"""
    return template("Hello <b>{{world}}</b>!")


if __name__ == "__main__":
    run(host="localhost", port=8080)
