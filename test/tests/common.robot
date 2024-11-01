*** Settings ***
Library    JSONLibrary
Library    Process
Library    KubeLibrary

*** Variables ***
${ROOT_URL}    http://84.3.93.248:30021
${TEST_NAMESPACE}    idk1


*** Keywords ***
Setup Suite
    ${ns_exist}=    Is Namespace Exists    ${TEST_NAMESPACE}
    Run Keyword If    ${ns_exist}    Delete Namespace    ${TEST_NAMESPACE}
    Create Namespace    ${TEST_NAMESPACE}
    ${ns_exist}=    Is Namespace Exists    ${TEST_NAMESPACE}
    Should Be True    ${ns_exist} 
    Helm Install Hfc    ${TEST_NAMESPACE}

Call Endpoint
    [Documentation]    Send a curl request with a JSON body. URL, METHOD, DATA,
    [Arguments]    ${URL_END}    ${method}    ${req_dat_path}
    ${full_url}=    Set Variable    ${ROOT_URL}${URL_END}
    ${command}=    Set Variable    curl -k -X ${method} ${full_url} -H "Content-Type: application/json" -H "Accept: application/json"
    ${command}=    Run Keyword If    '${req_dat_path}' != ''    Set Variable   ${command} -d @${req_dat_path}
    ${result}=    Run Process    ${command}    shell=True
    RETURN    ${result}

Set Up Hfc
    Get Healthcheck

Wait Until Pod Is Running
    [Arguments]    ${pod_name_pattern}    ${namespace}
    Wait Until Keyword Succeeds    2min    5sec 
    ...  Is Pod Running    ${POD_NAME_PATTERN}    ${NAMESPACE}

Is Pod Running
    [Arguments]    ${pod_name_pattern}    ${namespace}
    ${pods}=    List Namespaced Pod By Pattern    ${pod_name_pattern}    ${namespace}
    ${pod_name}=     Set Variable    ${pods[0].metadata.name}
    ${status}=    Read Namespaced Pod Status    ${pod_name}    ${NAMESPACE}
    Should Be True     '${status.phase}'=='Running'  

Is Namespace Exists
    [Arguments]    ${namespace}
    ${namespaces}=    Get Namespaces
    ${exists}=    Run Keyword And Return Status    Should Contain    ${namespaces}    ${namespace}
    RETURN    ${exists}

Create Namespace
    [Arguments]    ${namespace}
    ${result}=    Run Process    minikube kubectl -- create ns ${namespace}    shell=True
    Should Be Empty    ${result.stderr}

Delete Namespace
    [Arguments]    ${namespace}
    ${result}=    Run Process    minikube kubectl -- delete ns ${namespace}    shell=True
    Should Be Empty    ${result.stderr}

Helm Install Hfc
#helm install hfc ./hfc/helm -n test1
    [Arguments]    ${namespace}
    ${result}=    Run Process    helm install hfc ./hfc/helm -n ${namespace}    shell=True
    Should Be Empty    ${result.stderr}