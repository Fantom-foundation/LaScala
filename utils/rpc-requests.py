import yaml
import requests
import re
import logging
import sys
import getopt
import time

logging.basicConfig(format='%(asctime)s %(message)s', datefmt='%Y/%m/%d %H:%M:%S', level=logging.INFO)


def main(argv):
    with open('rpc-requests.yaml') as f:
        queries = yaml.safe_load(f)

    url = 'http://0.0.0.0:80'
    headers = {'Content-Type': 'application/json'}
    errors = 0
    count = 0

    # parse command line options:
    try:
        opts, args = getopt.getopt(argv, "rpc:", ["rpc-node-url="])
    except getopt.GetoptError:
        sys.exit(2)

    for opt, arg in opts:
        if opt in ("-rpc", "--rpc-node-url"):
            # use alternative algorithm:
            url = arg

    logging.info("Using URL: " + url)

    # Positional command line arguments (i.e. non optional ones) are
    # still available via 'args':
    logging.info("Positional args: " + args)

    for q in queries:
        time.sleep(100/1000)
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
    main(sys.argv[1:])
