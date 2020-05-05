from urllib.request import urlopen
from bs4 import BeautifulSoup
from pync import Notifier
import logging
import logging
import sys
import time
logger = logging.getLogger('mylogger')
logger.setLevel(logging.DEBUG)

formatter = logging.Formatter('[%(asctime)s] p%(process)s {%(pathname)s:%(lineno)d} %(levelname)s - %(message)s','%m-%d %H:%M:%S')

stdout=logging.StreamHandler(sys.stdout)
stdout.setFormatter(formatter)

logger.addHandler(stdout)

import signal

def sigint_handler(num, frame):
    print("received SIGINT!")
    sys.exit(0)


def sigtstp_handler(num, frame):
    print("received SIGTSTP!")
    sys.exit(0)


def sigterm_handler(num, frame):
    print("received SIGTERM!")
    sys.exit(0)


html_doc = """
<html><head><title>The Dormouse's story</title></head>
<body>
<p class="title"><b>The Dormouse's story</b></p>

<p class="story">Once upon a time there were three little sisters; and their names were
<a href="http://example.com/elsie" class="sister" id="link1">Elsie</a>,
<a href="http://example.com/lacie" class="sister" id="link2">Lacie</a> and
<a href="http://example.com/tillie" class="sister" id="link3">Tillie</a>;
and they lived at the bottom of a well.</p>

<p class="story">...</p>
"""
lastTitle=""
def crawler():
    global  lastTitle
    """
    :returns: @todo

    """
    # soup = BeautifulSoup(html_doc, 'html.parser')
    # print(soup.prettify())
    # print(soup.title)

    html = urlopen("http://www.gjgwy.org/guangxi/zkgg/sydw/")
    # webDoc = BeautifulSoup(html.read(),'html.parser')
    webDoc = BeautifulSoup(html,'html.parser')
    # webDocMaincon=webDoc.find_all("div", "main_con")
    # items=webDocMaincon.find_all("span","cnt")
    webDocMaincon=webDoc.find("div", "main_con")
    itemsTime=webDocMaincon.findAll("span","time2")
    itemsTitle=webDocMaincon.findAll("span","cnt")
    keywords=["钦州","灵山","南宁","北海","防城港"]
    matchList=[]
    newTitle=""
    for index,item in enumerate(itemsTime):
        rawWord=itemsTitle[index].text
        logger.info("index:%d %s %s"%(index,item.text,itemsTitle[index].text))
        if rawWord==lastTitle :
            logger.info("no new job apply is published for %s"%(keywords))
            break
        if index==0:
            lastTitle=rawWord
        
        for w in keywords:
            if w in rawWord:
                logger.info("index:%d %s %s [%s]"%(index,item.text,itemsTitle[index].text,itemsTitle[index].a['href']))
                # logger.info("itemsTitle[index].a:%s"%(itemsTitle[index].a['href']))
                one={}
                one["title"]=rawWord+"@"+item.text
                one["url"]=itemsTitle[index].a['href']
                matchList.append(one)
                break
                
    for item in matchList:
        logger.info("item:%s"%(item))
        Notifier.notify(item["title"], open=item["url"])
    # items=webDocMaincon.findAll("span")
    # for item in items:
    #     print(item.text)
    #     print(item['class'])

def ticker():
    t=30*60
    logger.info("ticker :%d seconds"%(t ))
    time.sleep(t)

if __name__ == '__main__':
    signal.signal(signal.SIGINT, sigint_handler)
    signal.signal(signal.SIGTSTP, sigtstp_handler)
    signal.signal(signal.SIGTERM, sigterm_handler)
    while True:
        crawler()
        ticker()

# main_con
# <div class="main_con">
#         <ul class="list">
#         <li><span class="cnt"><a href="http://www.gjgwy.org/201905/419594.html" target="_blank" style="">2019年广西特岗教师招聘8254人公告</a></span><span class="time2">2019-05-24</span></li>
