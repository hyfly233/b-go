# script.py
import sys

def main():
    print("Hello from Python!")
    if len(sys.argv) > 1:
        print("Arguments:", sys.argv[1:])

if __name__ == "__main__":
    main()
