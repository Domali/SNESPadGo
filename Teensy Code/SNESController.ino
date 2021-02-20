#define CHECK_BIT(var,pos) ((var) & (1<<(pos)))
int ret; //  32 bits of data to return

//All the pins for reading the SNES controller.  These are all digital read ports.
static int dataPin = 6;//red
static int latchPin = 4;//yellow
static int clockPin = 2;//blue

//The state that the color is currently in
int state;

void setup() {
   Serial.begin(57600);
   pinMode(dataPin, INPUT_PULLUP);
   pinMode(latchPin, INPUT);
   pinMode(clockPin, INPUT);
}

void loop(){
     ret = 0;
     //Wait for the latch to rise and fall to know that data is coming
     while(digitalReadFast(latchPin) == LOW){}//Wait for the latch to rise
     while(digitalReadFast(latchPin) == HIGH){}//Watch for the latch to fall
     //Wait for the clock to fall so we know when to read data off of the data pin.  Rinse and repeat 16 times.
     for(int i = 0;i<16;i++){
       while(digitalReadFast(clockPin)==HIGH){
       }
       ret |= digitalReadFast(dataPin) << i;
       while(digitalReadFast(clockPin)==LOW){
       }
     }
     Serial.println(ret,BIN);
}
