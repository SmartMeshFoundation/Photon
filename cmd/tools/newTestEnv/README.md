# newTestEnv

newTestEnv is  a tool to create a new environment for test.

1. deploy contract
2. register a new token
3. transfer money to A,B,C,D,E,F
4. create channels
    4.1 create channel A-B and save 100 both
    4.2 create channel B-C and save 50 both
    4.3 create channel C-E and save 100 both
    4.4 create channel A-D and save 100 both
    4.5 create channel B-D and save 100 both
    4.6 create channel D-F and save 100 both
    4.7 create channel F-E and save 100 both
    4.8 create channel C-F and save 50 both
    4.9     D-E 100

5. test path
5.1 A-B
5.2 A-B-C
5.3 A-B-A-D-E-C  A to c
5.4 A-B-C-B-D-E-F A to F