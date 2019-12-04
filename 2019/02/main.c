#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define BUFFERSIZE 1024
#define MAX_SIZE 4096
#define MAX_PROGRAM_SIZE 2048


int main(int argc, char** argv) {

    char buffer[BUFFERSIZE] = {0};
    char text[MAX_SIZE] = {0};

    int noun = -1;
    int verb = -1;

    int wanted_output = -1;

    FILE* input = NULL;
    if (argc == 0)
    {
        input = stdin;
    }
    else
    {
        if (argc > 1)
        {
            if (strcmp(argv[1], "-") == 0)
                input = stdin;
            else
                input = fopen(argv[1], "r");
        }

        if (argc == 3)
        {
            wanted_output = atoi(argv[2]);
        }

        if (argc == 4)
        {
            noun = atoi(argv[2]);
            verb = atoi(argv[3]);
        }
    }

    while(fgets(buffer, BUFFERSIZE, input))
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

    int min_noun = 0;
    int max_noun = 99;
    int found_noun = noun;

    int min_verb = 0;
    int max_verb = 99;
    int found_verb = verb;

    // reset noun and verb to random values if we have to find them
    if ((noun < 0) && (verb < 0) && (wanted_output > 0))
    {
        noun = 12;
        verb = 2;
    }

    while(1)
    {
        memcpy(memcode, program, sizeof(int)*MAX_PROGRAM_SIZE);

        if ((noun >= 0) && (verb >= 0))
        {
            fprintf(stderr, "Patching opcodes with noun=%d and verb=%d\n", noun, verb);
            memcode[1] = noun;
            memcode[2] = verb;
        }

        for(int ip=0; ip<pc; ip++)
        {
            int opc = memcode[ip];
            if (opc == 1 || opc == 2) {
                int a = memcode[++ip];
                int b = memcode[++ip];
                int o = memcode[++ip];

                if (opc == 1)
                {
                    fprintf(stderr, "[%d] = [%d] + [%d]\n", o, a, b);
                    memcode[o] = memcode[a] + memcode[b];
                }
                else
                {
                    fprintf(stderr, "[%d] = [%d] * [%d]\n", o, a, b);
                    memcode[o] = memcode[a] * memcode[b];
                }
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

        if (
            (wanted_output < 0)                     // we just want to run the program as is
            || (memcode[0] == wanted_output)        // we are looking for a specific output
            || (found_noun > 0 && found_verb > 0))  // we couldn't find an exact answer and this is our best guess
        {
            printf("output=%d, noun=%d, verb=%d, test=%d\n", memcode[0], memcode[1], memcode[2], 100*memcode[1]+memcode[2]);
            break;
        }
        else if (memcode[0] > wanted_output)
        {
            if (found_noun < 0)
            {
                fprintf(stderr, "output %d is greater, changing noun\n", memcode[0]);
                max_noun = noun;
                noun = (max_noun + min_noun) / 2;
                if (noun == min_noun)
                {
                    found_noun = noun;
                }
            }
            else if (found_verb < 0)
            {
                fprintf(stderr, "output %d is greater, changing verb\n", memcode[0]);
                max_verb = verb;
                verb = (max_verb + min_verb) / 2;
                if (verb == min_verb)
                {
                    found_verb = verb;
                }
            }
        }
        else if (memcode[0] < wanted_output)
        {
            if (found_noun < 0)
            {
                fprintf(stderr, "output %d is lower, changing noun\n", memcode[0]);
                min_noun = noun;
                noun = (max_noun + min_noun) / 2;
                if (noun == min_noun)
                {
                    found_noun = noun;
                }
            }
            else if (found_verb < 0)
            {
                fprintf(stderr, "output %d is lower, changing verb\n", memcode[0]);
                min_verb = verb;
                verb = (max_verb + min_verb) / 2;
                if (verb == min_verb) {
                    found_verb = verb;
                }
            }
        }
    };

    return 0;
}
