*** Settings ***
Resource    common.robot

Suite Setup    Setup Suite
Suite Teardown    Teardown Suite

Library    KubeLibrary
Library    Collections
Library    DateTime

*** Variables ***
${URL_END_REG_PRODUCT}    /prod
${URL_END_GET_PRODUCT}    /prod/8713600286072
${URL_END_GET_PRODUCT_NOT_EXISTING}    /prod/12340000
${URL_END_GET_CAT_LIST}    /prod/categories

${CAT_ABCDEF_JSON}    test/resources/categories/rerouting/ABCDEF.json
${CAT_ABG_JSON}    test/resources/categories/rerouting/ABG.json
${CAT_AXCDEF_JSON}    test/resources/categories/rerouting/AXCDEF.json
${CAT_AYZWEF_JSON}    test/resources/categories/rerouting/AYZWEF.json
${BASE_STATE_REROUTING}    test/resources/categories/rerouting/base_state.json
${REROUTED1X}    test/resources/categories/rerouting/rerouted1x.json
${REROUTED2X}    test/resources/categories/rerouting/rerouted2x.json

${REQ_JSON_NEW_PROD}    test/resources/new_prod_post_get/new_prod.json
${EXPECTED_JSON_NEW_PROD}    test/resources/new_prod_post_get/new_prod_expected.json

*** Test Cases ***
Category List When Empty
    ${result}=    Call Endpoint Filepath Data    ${URL_END_GET_CAT_LIST}    "GET"    ""
    Should Be Equal    ${result.stdout}    null

Category Rerouting Test
    Call Endpoint Filepath Data    ${URL_END_REG_PRODUCT}    "POST"    ${CAT_ABCDEF_JSON}
    Call Endpoint Filepath Data    ${URL_END_REG_PRODUCT}    "POST"    ${CAT_ABG_JSON}
    ${result}=    Call Endpoint Filepath Data    ${URL_END_GET_CAT_LIST}    "GET"    ""
    ${result_json}=    Convert String To Json    ${result.stdout}
    ${expected_json}=    Load Json From File    ${BASE_STATE_REROUTING}
    Lists Should Be Equal    ${result_json}    ${expected_json}    ignore_order=true

    # Reroute 1
    Call Endpoint Filepath Data    ${URL_END_REG_PRODUCT}    "POST"    ${CAT_AXCDEF_JSON}
    ${result}=    Call Endpoint Filepath Data    ${URL_END_GET_CAT_LIST}    "GET"    ""
    ${result_json}=    Convert String To Json    ${result.stdout}
    ${expected_json}=    Load Json From File    ${REROUTED1X}
    Lists Should Be Equal    ${expected_json}    ${result_json}    ignore_order=true
    
    # Reroute 2
    Call Endpoint Filepath Data    ${URL_END_REG_PRODUCT}    "POST"    ${CAT_AYZWEF_JSON}
    ${result}=    Call Endpoint Filepath Data    ${URL_END_GET_CAT_LIST}    "GET"    ""
    ${result_json}=    Convert String To Json    ${result.stdout}
    ${expected_json}=    Load Json From File    ${REROUTED2X}
    Lists Should Be Equal    ${expected_json}    ${result_json}    ignore_order=true

New Product Registering
    [Documentation]    Registers a new product then checks by getting the product by barcode
    ${result}=    Call Endpoint Filepath Data    ${URL_END_REG_PRODUCT}    "POST"    ${REQ_JSON_NEW_PROD}
    ${result_json}=    Convert String To Json    ${result.stdout}
    ${expected_json}=    Load Json From File    ${EXPECTED_JSON_NEW_PROD}
    Should Be Equal    ${result_json}    ${expected_json}
    
Get Product Which Was Registered
    ${got_product}=    Call Endpoint Filepath Data    ${URL_END_GET_PRODUCT}    "GET"    ""
    ${got_product_json}=    Convert String To Json    ${got_product.stdout}
    ${expected_json}=    Load Json From File    ${EXPECTED_JSON_NEW_PROD}
    Should Be Equal    ${got_product_json}    ${expected_json}

Get Product Which Was Not Registered
    ${got_product}=    Call Endpoint Filepath Data    ${URL_END_GET_PRODUCT_NOT_EXISTING}    "GET"    ""
    ${got_product_json}=    Convert String To Json    ${got_product.stdout}
    Should Be Empty    ${got_product_json}