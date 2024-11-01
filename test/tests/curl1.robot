*** Settings ***
Resource    common.robot

Suite Setup    Setup Suite
Suite Teardown    Teardown Suite

Library     KubeLibrary

*** Variables ***
${URL_END_REG_PRODUCT}    /prod
${URL_END_GET_PRODUCT}    /prod/8713600286072

${REQ_JSON}    test/resources/prod_post/req1.json
${EXPECTED_JSON}    test/resources/prod_post/req1_expected.json

*** Test Cases ***
New Product Registering
    [Documentation]    Registers a new product then checks by getting the product by barcode
    ${result}=    Call Endpoint    ${URL_END_REG_PRODUCT}    "POST"    ${REQ_JSON}
    ${result_json}=    Convert String To Json    ${result.stdout}
    ${expected_json}=    Load Json From File    ${EXPECTED_JSON}
    Should Be Equal    ${result_json}    ${expected_json}
    ${got_product}=    Call Endpoint    ${URL_END_GET_PRODUCT}    "GET"    ""
    ${got_product_json}=    Convert String To Json    ${got_product.stdout}
    Should Be Equal    ${got_product_json}    ${expected_json}

*** Keywords ***



Teardown Suite
    RETURN    ""