from bs4 import BeautifulSoup
from mitmproxy import ctx

# load in the javascript to inject
with open('bundle.js', 'r') as f:
    bundle_js = f.read()

def response(flow):
    if not ("text/html" in flow.response.headers["content-type"]) or not flow.response.status_code == 200: # only process 200 responses of html content
        return
    # inject the script tag
    html = BeautifulSoup(flow.response.text, 'lxml')
    container = html.head or html.body
    if container:
        script = html.new_tag('script', type='text/javascript')
        script.string = bundle_js
        container.insert(0, script)
        flow.response.text = str(html)
        ctx.log.info('Successfully injected the content.js in html.')


def request(flow):
    # print("-"*50 + "request headers:")
    # for k, v in flow.request.headers.items():
    #     print("%-20s: %s" % (k.upper(), v))
    flow.request.headers["user-agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"
    flow.request.headers["accept-language"] = "en-US,en;q=0.9"

