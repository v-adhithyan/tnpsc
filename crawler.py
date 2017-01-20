from selenium import webdriver
import os
import time
import sys
import argparse
import xml.etree.cElementTree as ET

def write(contents, filename):
    content = "\n".join(contents)
    with open("{}.txt".format(filename), "w") as f:
        f.write(content)

# Parse arguments
parser = argparse.ArgumentParser(description = "TNPSC Preparation Helper tool")
parser.add_argument('--url', help='URL to scrape', required=True)
parser.add_argument('--out', help="Output files", required=True)
parser.add_argument('--key', help="Key to scan in links", required=True)
args = parser.parse_args()

# Setup driver
chromedriver = "{}/driver/chromedriver".format(os.getcwd())
os.environ["webdriver.chrome.driver"] = chromedriver
driver = webdriver.Chrome(executable_path=chromedriver)

visited = []
contents = []

try:

    driver.get(args.url)
    visited.append(driver.current_url)
    i = 1
    while True:
        print "Page {}".format(i)
        i = i + 1

        links = driver.find_elements_by_xpath("//*[@href]")
        next_page = ""
        urls = []
        for link in links:
            url = link.get_attribute("href")
            text = link.text

            if text is not None and text.lower().find("next") >=0 :
                next_page = url

            if url.lower().find(args.key) >= 0:
                print url
                urls.append(url)
        driver.get(next_page)

        """try:
            #next_page = driver.find_elements_by_xpath("//li[@class='pager-next']")[0].get_attribute("href")
            #next_page = driver.find_element_by_link_text("next").get_attribute("href")
            next_page = driver.find_element_by_class_name("pager-next").get_attribute("href");
            print next_page
            driver.get(next_page)
        except:
            raise
            break"""

    driver.quit()

    write(urls, args.out)
except:
    driver.quit()
    raise
