// Basic demo for accelerometer readings from Adafruit MSA301

#include <Wire.h>
#include <Adafruit_MSA301.h>
#include <Adafruit_Sensor.h>
#include <HTTPClient.h>
#include <WiFiClientSecure.h>

#include "config.h"

Adafruit_MSA301 msa;
TwoWire MSATW = TwoWire(0);
WiFiClientSecure *client = new WiFiClientSecure;
// Add a scoping block for HTTPClient https to make sure it is destroyed before WiFiClientSecure *client is
HTTPClient https;


/*
  Function to set up the connection to the WiFi AP
*/
void setupWiFi() {
  Serial.print("Connecting to '");
  Serial.print(WIFI_SSID);
  Serial.print("' ...");

  WiFi.mode(WIFI_STA);
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("connected");

  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());

  randomSeed(micros());
}

void setup(void) {
  MSATW.begin(15, 13, 100000);
  Serial.begin(115200);
  while (!Serial) delay(10);     // will pause Zero, Leonardo, etc until serial console opens

  Serial.println("Starting fidgobject");

  // Try to initialize!
  if (! msa.begin(MSA301_I2CADDR_DEFAULT, &MSATW)) {
    Serial.println("Failed to find MSA301 chip");
    while (1) {
      delay(10);
    }
  }
  Serial.println("MSA301 Found and connected");
  msa.setDataRate(MSA301_DATARATE_1_HZ);

  client -> setCACert(ROOT_CA);
  https.setReuse(true);


}


/*
   Function to submit metrics to hosted Graphite
*/
void submitSensors(float x, float y, float z) {
  // build hosted metrics json payload
  String body = String("{") +
                "\"id\": \"" + ID + "\"," +
                "\"x\": \"" + x + "\"," +
                "\"y\": \"" + y + "\"," +
                "\"z\": \"" + z + "\"" +
                "}";

  Serial.println(body);

  
  if (client) {
      // submit POST request via HTTP
      https.begin(String("https://") + TIMEFIDGET_HOST + "/push");
      https.addHeader("Content-Type", "application/json");

      int httpCode = https.POST(body);
      if (httpCode > 0) {
        Serial.printf("timefidget [HTTP] POST...  Code: %d  Response: ", httpCode);
        https.writeToStream(&Serial);
        Serial.println();
      } else {
        Serial.printf("timefidget [HTTP] POST... Error: %s\n", https.errorToString(httpCode).c_str());
      }

      https.end();

  } else {
    Serial.println("Unable to create client");
  }

}



void loop() {
  // reconnect to WiFi if required
  if (WiFi.status() != WL_CONNECTED) {
    WiFi.disconnect();
    yield();
    setupWiFi();
  }


  /* Or....get a new sensor event, normalized */
  sensors_event_t event;
  msa.getEvent(&event);

  /* Display the results (acceleration is measured in m/s^2) */
  //  Serial.print("\t\tX: "); Serial.print(event.acceleration.x);
  //  Serial.print(" \tY: "); Serial.print(event.acceleration.y);
  //  Serial.print(" \tZ: "); Serial.print(event.acceleration.z);
  //  Serial.println(" m/s^2 ");

  Serial.println();

  submitSensors(event.acceleration.x, event.acceleration.y, event.acceleration.z);

  delay(5000);
}
