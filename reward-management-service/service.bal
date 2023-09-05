import ballerina/log;
import ballerina/http;




type RewardDeal record {|
    string userId;
    string dealId;
    boolean acceptedTC;
|};


# A service representing a network-accessible API
# bound to port `9090`.
service / on new http:Listener(9090) {

    resource function post create(@http:Payload RewardDeal payload) returns error?  {
        log:printInfo("RewardDeal: ", payload = payload);

    }
}
