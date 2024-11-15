*** Settings ***
Library    JSONLibrary
Library    Process
Library    KubeLibrary

*** Variables ***
${HFC_PORT}    30101
${ROOT_URL}    http://192.168.49.2:${HFC_PORT}
#${ROOT_URL}    http://192.168.0.69:1323
${HFC_PODNAME}    hfc

${TEST_NAMESPACE}    test
${_mansetup}    false

${MONGO_PODNAME}    mongo
${MONGO_HOST}    mongo-mongodb.${TEST_NAMESPACE}.svc.cluster.local
${MONGO_PORT}    27017
${MONGO_PWD}    test1234

${OFF_IMPORT_PROD_SUSHI}    test/resources/mongo_off_default/sushi.json
${OFF_IMPORT_PROD_ZOTT}    test/resources/mongo_off_default/zott.json

${MONGO_TESTPOD_NAME}    mongodb-test-pod

*** Keywords ***
Call Endpoint Direct Data
    [Documentation]    Send a curl request with a JSON body from parameter. URL, METHOD, DATA,
    [Arguments]    ${URL_END}    ${method}    ${req_dat_path}
    ${full_url}=    Set Variable    ${ROOT_URL}${URL_END}
    ${command}=    Set Variable    curl -k -X ${method} ${full_url} -H "Content-Type: application/json" -H "Accept: application/json" -d '${req_dat_path}'
    Log To Console    \r\nUsing command: ${command}
    ${result}=    Run Process    ${command}    shell=True
    RETURN    ${result}

Call Endpoint Filepath Data
    [Documentation]    Send a curl request with a JSON body. URL, METHOD, DATA from file,
    [Arguments]    ${URL_END}    ${method}    ${req_dat_path}
    ${full_url}=    Set Variable    ${ROOT_URL}${URL_END}
    ${command}=    Set Variable    curl -k -X ${method} ${full_url} -H "Content-Type: application/json" -H "Accept: application/json"
    ${command}=    Run Keyword If    '${req_dat_path}' != ''    Set Variable   ${command} -d @${req_dat_path}
    Log To Console    \r\nUsing command: ${command}
    ${result}=    Run Process    ${command}    shell=True
    RETURN    ${result}

Setup Suite
    ${ns_exist}=    Is Namespace Exists    ${TEST_NAMESPACE}
    Run Keyword If    ${ns_exist}    Delete Namespace    ${TEST_NAMESPACE}
    Create Namespace    ${TEST_NAMESPACE}
    ${ns_exist}=    Is Namespace Exists    ${TEST_NAMESPACE}
    Should Be True    ${ns_exist} 
    Install Mongodb    ${TEST_NAMESPACE}
    Install Mongo Test Pod    ${TEST_NAMESPACE}
    Copy File To TestPod    ${OFF_IMPORT_PROD_SUSHI}    ${TEST_NAMESPACE}    ${MONGO_TESTPOD_NAME}
    TestPod Import Off Product    ${TEST_NAMESPACE}    ${MONGO_TESTPOD_NAME}
    Copy File To TestPod    ${OFF_IMPORT_PROD_ZOTT}    ${TEST_NAMESPACE}    ${MONGO_TESTPOD_NAME}
    TestPod Import Off Product    ${TEST_NAMESPACE}    ${MONGO_TESTPOD_NAME}
    Helm Install Hfc    ${TEST_NAMESPACE}

Teardown Suite
    ${ns_exist}=    Is Namespace Exists    ${TEST_NAMESPACE}
    Run Keyword If    ${ns_exist}    Delete Namespace    ${TEST_NAMESPACE}

Wait Until Pod Is Running And Ready
    [Arguments]    ${pod_name_pattern}    ${namespace}
    Log To Console    Waiting for pod: ${pod_name_pattern}
    Wait Until Keyword Succeeds    1min    5sec 
    ...  Is Pod Phase Running    ${POD_NAME_PATTERN}    ${NAMESPACE}
    Wait Until Keyword Succeeds    2min    5sec 
    ...  Is Pod Containers Ready    ${POD_NAME_PATTERN}    ${NAMESPACE}

Is Pod Phase Running
    [Arguments]    ${pod_name_pattern}    ${namespace}
    ${pods}=    List Namespaced Pod By Pattern    ${pod_name_pattern}    ${namespace}
    ${pod_name}=     Set Variable    ${pods[0].metadata.name}
    ${status}=    Read Namespaced Pod Status    ${pod_name}    ${NAMESPACE}
    Log To Console    waiting for pod running phase
    Should Be True     '${status.phase}'=='Running'

Is Pod Containers Ready
    [Arguments]    ${pod_name_pattern}    ${namespace}
    ${pods}=    List Namespaced Pod By Pattern    ${pod_name_pattern}    ${namespace}
    ${containers_status}=    Filter Pods Containers Statuses By Name    ${pods}    .*
    Log To Console    waiting for pod readiness x/x
    FOR    ${container}    IN    @{containers_status}
        Should Be True    ${container.ready}
    END

Is Namespace Exists
    [Arguments]    ${namespace}
    ${namespaces}=    Get Namespaces
    ${exists}=    Run Keyword And Return Status    Should Contain    ${namespaces}    ${namespace}
    RETURN    ${exists}

Create Namespace
    [Arguments]    ${namespace}
    Log To Console    Creating namespace ${namespace}
    ${result}=    Run Process    minikube kubectl -- create ns ${namespace}    shell=True
    Should Be Empty    ${result.stderr}

Delete Namespace
    [Arguments]    ${namespace}
    Log To Console    Deleting namespace ${namespace}
    ${result}=    Run Process    minikube kubectl -- delete ns ${namespace}    shell=True
    Should Be Empty    ${result.stderr}
    
Helm Install Hfc
    [Arguments]    ${namespace}
    Log To Console    Installing hfc chart
    ${command}=    Set Variable    helm install hfc ./hfc/helm -n ${namespace} --set env.mongo.host\=${MONGO_HOST},env.mongo.port\=${MONGO_PORT},service.port\=${HFC_PORT},service.nodePort\=${HFC_PORT}
    Log To Console    Running cmd: ${command}
    ${result}=    Run Process    ${command}    shell=True
    Wait Until Pod Is Running And Ready    ${HFC_PODNAME}    ${TEST_NAMESPACE}

Install Mongodb
    [Arguments]    ${namespace}
    Log To Console    Installing mongodb
    ${result}=    Run Process    bash ./mongo/myStart.sh ${namespace} ${MONGO_PWD}    shell=True
    Wait Until Pod Is Running And Ready    ${MONGO_PODNAME}    ${TEST_NAMESPACE}

Install Mongo Test Pod
    [Arguments]    ${namespace}
    Log To Console    Installing mongo testpod
    ${command}=    Set Variable    minikube kubectl -- apply -f mongo/manual-test-pod.yaml -n ${namespace}
    Log To Console    Running cmd: ${command}
    ${result}=    Run Process    ${command}    shell=True
    Wait Until Pod Is Running And Ready    mongodb-test-pod    ${namespace}

Copy File To TestPod
    [Arguments]    ${filePathToCopy}    ${namespace}    ${podname}
    Log To Console    Copying file to mongot testpod
    ${command}=    Set Variable    minikube kubectl -- cp ${filePathToCopy} ${namespace}/${podname}:/etc/mongoimport.json
    Log To Console    Running cmd: ${command}
    ${result}=    Run Process    ${command}    shell=True

TestPod Import Off Product
    [Arguments]    ${namespace}    ${podname}
    Log To Console    Mongo import from json
    ${command}=    Set Variable    minikube kubectl -- exec pod/${podname} -n ${namespace} -- mongoimport -u root -p ${MONGO_PWD} --authenticationDatabase admin --host ${MONGO_HOST} --port ${MONGO_PORT} --db off --collection products --file /etc/mongoimport.json
    Log To Console    Running cmd: ${command}
    ${result}=    Run Process    ${command}    shell=True

*** Test Cases ***
Actually Setup
    [Documentation]     Using robot to set up env for manual purposes
    IF    '${_mansetup} == true'
            Setup Suite
    END