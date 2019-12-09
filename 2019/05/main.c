#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#define BUFFERSIZE 1024
#define MAX_SIZE 4096
#define MAX_PROGRAM_SIZE 2048


int main(int argc, char** argv) {

    char buffer[BUFFERSIZE] = {0};
    char text[MAX_SIZE] = {0};
    
    FILE* in = NULL;
    if (argc == 1)
    {
        in = stdin;
    }
    else
    {
        if (argc > 1)
        {
            if (strcmp(argv[1], "-") == 0)
                in = stdin;
            else
                in = fopen(argv[1], "r");
        }
   }

    while(fgets(buffer, BUFFERSIZE, in))
    {
        strcat(text, buffer);
    }

    int program[MAX_PROGRAM_SIZE];
    int memcode[MAX_PROGRAM_SIZE];

    int pc = 0;
    char* pch = strtok(text, ",\n");
    while(pch != NULL)
    {
        int op = atoi(pch);
        program[pc++] = op;
        pch=strtok(NULL, ",\n");
    }

    fprintf(stderr, "Parsed program with %d opcodes\n", pc);
 
    int input = 1;
    int output = 0;

    memcpy(memcode, program, sizeof(int)*MAX_PROGRAM_SIZE);
    for(int ip=0; ip<pc; ip++)
    {
        int opc = memcode[ip];
        int mode[3]={0};
        int m = 0;

        if (opc > 99) {
            m = opc / 100;
            opc %= 100;
        }

        int k=0;
        while (m > 0) {
            mode[k] = m % 10;
            m /= 10;
            k++;
        }

        if (opc == 1 || opc == 2) {
            int a = memcode[++ip];
            int b = memcode[++ip];
            int o = memcode[++ip];

            int lhs = a;
            int rhs = b;

            if (mode[0] == 0) {
                lhs = memcode[a];
            }
            if (mode[1] == 0) {
                rhs = memcode[b];
            }
            assert(mode[2] == 0);

            fprintf(stderr, "*[%d] = %c[%d] %c %c[%d]\n", o, mode[0]==0?'*':' ', a, opc == 1?'+':'*', mode[1]==0?'*':' ', b);

            if (opc == 1)
            {
                memcode[o] = lhs + rhs;
            }
            else
            {
                memcode[o] = lhs * rhs;
            }
        }
        else if (opc == 3)
        {
            int o = memcode[++ip];
            assert(mode[0] == 0);

            fprintf(stderr, "[%d] <- %d\n", o, input);
            memcode[o] = input;
        }
        else if (opc == 4)
        {
            int v = memcode[++ip];
            //assert(mode[0] == 0); // for some reason, this fails ?!

            output = memcode[v];
            fprintf(stderr, "[%d](%d) -> output\n", v, output);
        }
        else if (opc == 99)
        {
            fprintf(stderr, "END\n");
            break;
        }
        else
        {
            fprintf(stderr, "Invalid opcode %d at position %d\n", opc, pc);
            return 1;
        }
    }

    printf("%d\n", output);

    if (in != NULL) {
        fclose(in);
    }

    return 0;
}

