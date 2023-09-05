import ballerina/http;
import ballerina/log;

# A service representing a network-accessible API
# bound to port `9090`.
service / on new http:Listener(9090) {

    resource function post .(@http:Payload RewardConfirmationEvent payload) returns error? {
        log:printInfo("reward confirmation received", rewardConfirmation = payload);
    }
}

public type RewardConfirmationEvent record {|
    string userId;
    string rewardId;
    string rewardConfirmationNumber;
|};
