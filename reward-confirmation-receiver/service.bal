import ballerina/http;
import ballerina/log;
import ballerinax/mysql;
import ballerina/sql;
import ballerinax/mysql.driver as _;

# A service representing a network-accessible API
# bound to port `9090`.
service / on new http:Listener(9090) {

    resource function post confirm(@http:Payload RewardConfirmationEvent payload) returns error? {
        log:printInfo("reward confirmation received", rewardConfirmation = payload);

        log:printInfo("generate qr code for: ", rewardConformationNumber = payload.rewardConfirmationNumber);

        http:Client qrCodeClient = check new (qrcodeApiEndoint, {
            auth: {
                tokenUrl: qrCodeTokenEndpoint,
                clientId: qrCodeClientId,
                clientSecret: qrCodeClientSecret
            }
        });

        http:Response httpResponse = check qrCodeClient->/qrcode(content = payload.rewardConfirmationNumber);

        // Get the byte payload from the response
        byte[] imageContent = check httpResponse.getBinaryPayload();

        // Insert the image content into the MySQL database using the existing mysqlEndpoint
        sql:ParameterizedQuery insertQuery = `INSERT INTO reward_confirmation (id, reward_id, user_id, reward_confirmation_qrcode) VALUES (0, ${payload.rewardId}, ${payload.userId}, ${imageContent})`;
        sql:ExecutionResult result = check mysqlEndpoint->execute(insertQuery);

        if (result.affectedRowCount > 0) {
            log:printInfo("image successfully saved to the user profile");
        } else {
            log:printError("failed to save the image to the user profile");
        }

    }
}

public type RewardConfirmationEvent record {|
    string userId;
    string rewardId;
    string rewardConfirmationNumber;
|};

configurable string qrcodeApiEndoint = ?;
configurable string qrCodeClientId = ?;
configurable string qrCodeClientSecret = ?;
configurable string qrCodeTokenEndpoint = ?;

configurable string dbhost = ?;
configurable string dbuser = ?;
configurable string dbpwd = ?;
configurable string database = ?;

public mysql:Client mysqlEndpoint = check new (host = dbhost, user = dbuser, password = dbpwd, database = database);
