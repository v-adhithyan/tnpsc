from selenium import webdriver
import os
import time
import sys
import argparse
import xml.etree.cElementTree as ET

def write_to_xml(contents, filename):
    resources = ET.Element("resources")
    string_array = ET.SubElement(resources, "string-array", name=filename)

    for content in contents:
        #ET.SubElement(string_array, "item").text = str(content).encode("utf-8").strip()
        ET.SubElement(string_array, "item").text = content

    tree = ET.ElementTree(resources)
    tree.write("{}.xml".format(filename))

# Parse arguments
parser = argparse.ArgumentParser(description = "TNPSC Preparation Helper tool")
parser.add_argument('--url', help='URL to scrape', required=True)
parser.add_argument('--out', help="Output xml file name without xml extension", required=True)
args = parser.parse_args()

# Setup driver
chromedriver = "{}/driver/chromedriver".format(os.getcwd())
os.environ["webdriver.chrome.driver"] = chromedriver
driver = webdriver.Chrome(executable_path=chromedriver)

visited = []
contents = []

try:

    driver.get(args.url)

    while driver.current_url not in visited:
        visited.append(driver.current_url)

        paras = driver.find_elements_by_xpath("//p[@style='line-height: 150%']")

        n = len(paras)
        point = 1
        for i,p in enumerate(paras):
            if i == n-1:
                try:
                    next_page = driver.find_element_by_link_text("Next").get_attribute("href")
                    driver.get(next_page)
                except:
                    break
            else:
                content = p.text

                if len(content) > 30:
                    contents.append(content)
                    point = point + 1
                    print content

    driver.quit()

    write_to_xml(contents, args.out)
except:
    driver.quit()
    raise
