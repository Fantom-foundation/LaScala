import yaml
import requests
import re
import logging

logging.basicConfig(format='%(asctime)s %(message)s', datefmt='%Y/%m/%d %H:%M:%S', level=logging.INFO)


def main():
    with open('rpc-requests.yaml') as f:
        queries = yaml.safe_load(f)

    url = 'http://0.0.0.0:80'
    headers = {'Content-Type': 'application/json'}
    errors = 0
    count = 0

    for q in queries:
        logging.info('Processing ' + q['name'])
        response = requests.post(url, data=q['body'], headers=headers)
        expected = re.compile(q['result'])
        if expected.search(response.text):
            logging.info("\tPassed")
        else:
            logging.info("\tFailed")
            logging.info('\t\trequest:  ' + q['body'])
            logging.info('\t\texpect:   ' + q['result'])
            logging.info('\t\tresponse: ' + response.text)
            errors = errors + 1
        count = count + 1

    logging.info(str(count - errors) + '/' + str(count) + ' tests passed.')
    if errors > 0:
        exit(1)


if __name__ == '__main__':
    main()
