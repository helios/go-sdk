import sys
import os
import json

PROXY_LIBS_FILE_NAME = 'proxy_libs_versions.json'
PROXY_LIBS_PREFIX = 'github.com/helios/go-sdk/proxy-libs'

def main():
    proxy_lib_name = sys.argv[1]
    proxy_lib_version = sys.argv[2]
    min_version_supported = sys.argv[3]
    with open(PROXY_LIBS_FILE_NAME) as f:
        proxy_libs_versions = json.load(f)
        proxy_libs_versions[("%s/%s" % (PROXY_LIBS_PREFIX,proxy_lib_name))][min_version_supported] = proxy_lib_version
    os.remove(PROXY_LIBS_FILE_NAME)
    with open(PROXY_LIBS_FILE_NAME, 'w') as output:
        json.dump(proxy_libs_versions, output)

if __name__ == "__main__":
    main()   
 