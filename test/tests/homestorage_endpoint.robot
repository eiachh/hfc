*** Settings ***
Resource    common.robot

Suite Setup    Setup Suite
Suite Teardown    Teardown Suite

Library    KubeLibrary
Library    Collections
Library    DateTime

*** Variables ***

${URL_END_GET_HS}    /hs
${URL_END_POST_HS_ITEM}    /hs/4014500513010
${URL_END_PUT_HS_ITEM}    /hs/4014500513010

${EXPECTED_JSON_GET_HS_EMPTY}    test/resources/homestorage/empty_hs_resp.json
${EXPECTED_JSON_POST_HS_PROD}    test/resources/homestorage/new_unregistered_food_2amnt.json
${EXPECTED_JSON_POST_HS_PROD_RESP}    test/resources/homestorage/hs_put_unregistered_resp.json


*** Test Cases ***
Delete Attempt Non Existing Item
    ${uuid}=    Set Variable    b45592fd-473d-4a5d-89fa-c1c0506ce91d
    ${dict}=    Create Dictionary    uuid=${uuid}
    ${dict_as_json_str}=    Convert Json To String    ${dict}
    ${result}=    Call Endpoint Direct Data    ${URL_END_PUT_HS_ITEM}    "DELETE"    ${dict_as_json_str}
    ${result_stdout}=    Set Variable    ${result.stdout}
    Should Be Equal As Strings    ${result_stdout}    \{\}

Get Empty Hs
    ${got_hs}=    Call Endpoint Filepath Data    ${URL_END_GET_HS}    "GET"    ""
    ${got_hs_json}=    Convert String To Json    ${got_hs.stdout}
    ${expected_json}=    Load Json From File    ${EXPECTED_JSON_GET_HS_EMPTY}
    Should Be Equal    ${got_hs_json}    ${expected_json}

Add To HS Unregistered Product
    ${post_hs_resp}=    Call Endpoint Filepath Data    ${URL_END_POST_HS_ITEM}    "POST"    ${EXPECTED_JSON_POST_HS_PROD}
    ${post_hs_resp_json}=    Convert String To Json    ${post_hs_resp.stdout}
    ${expected_json}=    Load Json From File    ${EXPECTED_JSON_POST_HS_PROD_RESP}
    Should Be Equal    ${post_hs_resp_json}    ${expected_json}

Get Hs After Adding Product
    ${got_hs}=    Call Endpoint Filepath Data    ${URL_END_GET_HS}    "GET"    ""
    ${got_hs_json}=    Convert String To Json    ${got_hs.stdout}
    ${hs_items}=    Get Value From Json    ${got_hs_json}    $.home_storage_items
    ${first_hs_item}=    Set Variable    ${hs_items[0]}
    ${zott}=    Get From Dictionary    ${first_hs_item}    4014500513010
    ${first_zott}=    Set Variable    ${zott[0]}
    ${zott_amount}=    Get Length    ${zott}
    Should Be Equal As Strings    2    ${zott_amount}
    Dictionary Should Contain Key    ${first_zott}    uuid
    Should Not Be Empty   ${first_zott}[uuid]
    Dictionary Should Contain Key    ${first_zott}    acquired
    Should Not Be Empty   ${first_zott}[acquired]
    Dictionary Should Contain Key    ${first_zott}    expires
    Should Not Be Empty   ${first_zott}[expires]

Modify Hs Item 
    ${got_hs}=    Call Endpoint Filepath Data    ${URL_END_GET_HS}    "GET"    ""
    ${got_hs_json}=    Convert String To Json    ${got_hs.stdout}
    ${hs_items}=    Get Value From Json    ${got_hs_json}    $.home_storage_items
    ${first_hs_item}=    Set Variable    ${hs_items[0]}
    ${zott}=    Get From Dictionary    ${first_hs_item}    4014500513010
    ${first_zott}=    Set Variable    ${zott[0]}
    
    ${uuid}=    Set Variable    ${first_zott}[uuid]
    ${acquired}=    Set Variable    ${first_zott}[acquired]
    ${changed_date}=    Get Current Date    result_format=%Y-%m-%dT%H:%M:%S.1Z
    
    ${dict}=    Create Dictionary    uuid=${uuid}    acquired=${acquired}    expires=${changed_date}
    ${dict_as_json_str}=    Convert Json To String    ${dict}
    ${result}=    Call Endpoint Direct Data    ${URL_END_PUT_HS_ITEM}    "PUT"    ${dict_as_json_str}
    
    ${got_hs_after}=    Call Endpoint Filepath Data    ${URL_END_GET_HS}    "GET"    ""
    ${got_hs_json_after}=    Convert String To Json    ${got_hs_after.stdout}
    ${hs_items_after}=    Get Value From Json    ${got_hs_json_after}    $.home_storage_items
    ${first_hs_item_after}=    Set Variable    ${hs_items_after[0]}
    ${zott_after}=    Get From Dictionary    ${first_hs_item_after}    4014500513010
    ${first_zott_after}=    Set Variable    ${zott_after[0]}

    Should Be Equal    ${zott_after[0]}[uuid]    ${uuid}
    Should Be Equal    ${zott_after[0]}[acquired]    ${acquired}
    Should Be Equal    ${zott_after[0]}[expires]    ${changed_date}

Delete Hs Item 
    #Get a hs item
    ${got_hs}=    Call Endpoint Filepath Data    ${URL_END_GET_HS}    "GET"    ""
    ${got_hs_json}=    Convert String To Json    ${got_hs.stdout}
    ${hs_items}=    Get Value From Json    ${got_hs_json}    $.home_storage_items
    ${first_hs_item}=    Set Variable    ${hs_items[0]}
    ${zott}=    Get From Dictionary    ${first_hs_item}    4014500513010
    ${first_zott}=    Set Variable    ${zott[0]}
    
    #Pre-check
    ${original_uuid}=    Set Variable    ${first_zott}[uuid]
    ${zott_amount}=    Get Length    ${zott}
    Should Be Equal As Strings    2    ${zott_amount}

    #Delete by dict with uuid
    ${dict}=    Create Dictionary    uuid=${original_uuid}
    ${dict_as_json_str}=    Convert Json To String    ${dict}
    Call Endpoint Direct Data    ${URL_END_PUT_HS_ITEM}    "DELETE"    ${dict_as_json_str}

    #Get hs items after delete
    ${got_hs_after}=    Call Endpoint Filepath Data    ${URL_END_GET_HS}    "GET"    ""
    ${got_hs_json_after}=    Convert String To Json    ${got_hs_after.stdout}
    ${hs_items_after}=    Get Value From Json    ${got_hs_json_after}    $.home_storage_items
    ${first_hs_item_after}=    Set Variable    ${hs_items_after[0]}
    ${zott_after}=    Get From Dictionary    ${first_hs_item_after}    4014500513010
    ${zott_amount_after}=    Get Length    ${zott_after}
    
    Should Be Equal As Strings    1    ${zott_amount_after}
    FOR    ${one_of_zott}    IN    @{zott_after}
        ${item_uuid}=    Get From Dictionary    ${one_of_zott}    uuid
        ${is_match}=    Evaluate    "${item_uuid}" == "${original_uuid}"
        Should Not Be Equal As Strings    ${is_match}    True
    END