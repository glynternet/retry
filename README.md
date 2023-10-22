# retry

The command for retrying commands.

Pass commands to `retry` with configurable backoff settings, allowing retry of operations whilst not overloading a system.

## Use
```
Usage: retry [--] command and args

Usage of retry:
      --initial-backoff duration   initial backoff duration. (default 1s)
      --max-attempts uint          upper limit of number of attempts. 0 indicates no limit.
      --max-backoff duration       upper limit of backoff duration. 0 indicates no limit.
      --multiplier float           multiplier to apply after each failed attempt. (default 2)
      --randomisation float        randomisation to apply to the multiplication of each backoff
  -v, --version                    print the version
```

## Example

_Append to /tmp/lines once and fail script if less than three lines_
```
$ cat ./run.sh 
#!/bin/bash
set -euf -o pipefail
date >> /tmp/lines
test "$(wc -l /tmp/lines)" = 3
```

_Double the backoff time each iteration_
```
$ retry --initial-backoff 1s --multiplier 2 ./run.sh && echo done
error:  exit status 1
error:  exit status 1
done
```

_Result_
```
$ cat /tmp/lines 
Sat 21 Oct 21:17:37 MDT 2023
Sat 21 Oct 21:17:38 MDT 2023
Sat 21 Oct 21:17:40 MDT 2023
```
