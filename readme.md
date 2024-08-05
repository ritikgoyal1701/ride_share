Ride-Share Application

<b>Pre-requisites</b>

Before you start, ensure you have the following installed on your system:
   1. Go-lang programming language is installed in your system.
   2. Mongo db is installed.

Additionally, you need to download the following dependencies:
   1. github.com/go-redis/redis/v8
   2. github.com/google/uuid
   3. github.com/golang-jwt/jwt
   4. go.mongodb.org/mongo-driver/bson/primitive

<b>Running the Program</b>

To run the program, use the following command:
<i>go run *.go</i>

<b>Assumptions</b>

Following assumptions were made, which can be improved in the future.
1. One account is login in 1 device only.

<b>Enhancements</b>

To improve the application, consider implementing the following enhancements:
1. Caching: Implement caching to improve response times.
2. Real-time Updates: Use Pusher or Webhooks for a more enhanced user experience.
3. Asynchronous Processing: Implement workers using AWS SQS for decoupled tasks.
4. Payments: Add various payment functionalities.
5. Configuration Management: Create a configuration file to manage settings currently hard-coded in the program, such as:
   1. MongoDB's configuration
   2. Pricing
   3. Surge charges
   4. Redis configuration
6. User Experience Improvements: Enhance the user experience with functionalities like:
   1. Password reset
   2. Email receipts
   3. GST invoicing support