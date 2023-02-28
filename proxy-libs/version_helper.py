import sys
import os
import json

PROXY_LIBS_FILE_NAME = 'proxy_libs_versions.json'
PROXY_LIBS_PREFIX = 'github.com/helios/go-sdk/proxy-libs'

def main():
    proxy_lib_name = sys.argv[1]
    proxy_lib_version = sys.argv[2]
    min_supported_version = sys.argv[3]
    updated = False
    with open(PROXY_LIBS_FILE_NAME) as f:
        proxy_libs_versions = json.load(f)
        proxy_lib_full_path = ("%s/%s" % (PROXY_LIBS_PREFIX,proxy_lib_name))
        min_supported_versions = proxy_libs_versions[proxy_lib_full_path]
        for version in min_supported_versions:
            if version["minSupportedVersion"] == min_supported_version:
                version["proxyLibVersion"] = proxy_lib_version
                updated = True

    os.remove(PROXY_LIBS_FILE_NAME)
    if not updated:
        sys.exit("didn't find relevant tag to change")
    with open(PROXY_LIBS_FILE_NAME, 'w') as output:
        json.dump(proxy_libs_versions, output)

if __name__ == "__main__":
    main()  
    