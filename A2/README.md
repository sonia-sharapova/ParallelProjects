# Homework \#2

**See gradescope for due date**

This homework is intended to serve as an introduction to using atomic
operations and importance of understanding the architecture of the
parallel hardware a parallel application is running on.

### Programming Questions

For this assignment you are **only** allowed to use the following Go
concurrent constructs:

  - `go` statement
  - `sync/atomic` package. You may use any of the atomic operations.
  - `sync.WaitGroup` and its associated methods.

As the course progresses, you will be able to use other constructs
provided by Go such as channels, conditional variables, mutexes, etc. I
want you to learn the basics and build up from there. This way you'll
have a good picture of how parallel programs and their constructs are
built from each other. **If you are unsure about whether you are able to
use a language feature then please ask on Slack or during office hours
before using it**.

## Static Task Decomposition

For problems 1 and 2, you will be implementing your first parallel
programs in the course. Each goroutine will need to perform a certain
amount of work to complete. *But how should the work be distributed to
each goroutine*? For this assignment, we use the task distribution
technique we saw during class, which equally distributes the amount of
work to each goroutine.

## Problem 0: Set Up Cluster Access

I encourage you to set up your access to the Linux cluster, if you have not
done so already for other classes. Make sure you know your credentials,
know how to log in via `ssh`, and try to run a simple Go program, for example
the week 2 examples. Knowing that this all works will help you immensely with
the projects.

## Problem 1: Estimating Pi using the Leibniz formula

Inside the `pi` directory, open the file called `pi.go` and write a
concurrent Go program that uses an infinite series method to estimate pi.
One of the most widely known infinite series that can be used for this is
the Leibniz formula, as described here:

  - [How To Make
    Pi](https://en.wikipedia.org/wiki/Leibniz_formula_for_π)

Your program should use this formula, and must use the following usage and
have the these required command-line arguments

``` go
const usage = "Usage: pi interval threads\n" +
"    interval = the number of iterations to perform\n" +
"    threads = the number of threads (i.e., goroutines to spawn)\n"
```

The main goroutine should read in the `interval` and `threads` argument
and **only** print the estimate of pi. The program **must** use the
fork-join pattern described in class and a static task distribution
technique. The threshold for running the sequential version is 1000
intervals (even though the threads value is always set); otherwise, you
must run the parallel implementation.

**Assumptions**: You can assume that `interval` and `threads` will
always be given as integers and be provided. Thus, you do not need to
perform error checking on those arguments. You also will not need to
print out the usage statement for this problem. The usage is given here
for clarification purposes.

Sample Runs ($: is just mimicking the command line) :

    $: go run hw2/pi 100 2
    3.1315929036
    $: go run hw2/pi 10 1
    3.0418396189
    $: go run hw2/pi 10 3
    3.0418396189

The tests for this problem is inside `pi_test.go`. Please go back to
homework \#1 if you are having trouble running tests.

## Problem 2: Waitgroup Implementation

Read the following about go interfaces before starting this problem:

  - <https://gobyexample.com/interfaces>
  - <https://www.alexedwards.net/blog/interfaces-explained>

Inside the `waitgroup` directory, open the file called `waitgroup.go`.
Consider the following interface and function:

``` go
type WaitGroup interface {
  Add(amount uint)
  Done()
  Wait()
}

// NewWaitGroup returns a instance of a waitgroup
// This instance must be a pointer and should not
// be copied after creation.
func NewWaitGroup() WaitGroup {
  //IMPLEMENT ME!
}
```

Provide an implementation of a `WaitGroup` that mimics the Go
implementation:

  - <https://golang.org/pkg/sync/#WaitGroup>

The only difference is that our implementation requires uint (rather
than just int) to be passed to the `Add` method. Remember, you can only
use atomic operations in this implementation. Don't overthink this
problem. Many of these methods can be implemented in very few lines.
Note that there are various types of int (64 or 32 bit, signed or
unsigned), and there are special atomic operations for each type. The
test does not check which type you use inside your waitgroup
implementation.

The tests for this problem is inside `waitgroup_test.go`. Please go back
to homework \#1 if you are having trouble running tests.


### Program Performance

Do not worry about the performance of your programs for this week. We
will reexamine these implementations for the next homework assignment
once you a better understanding about parallel performance, which is
covered in module 3.

### Grading

Programming assignments will be graded according to a general rubric.
Specifically, we will assign points for completeness, correctness,
design, and style. (For more details on the categories, see Canvas.

The exact weights for each category will vary from one assignment to
another. For this assignment, the weights will be:

> **Note**
> Your code **must** be deterministic. This note means that if your code
> is running and passing the tests on your local machine or on CS remote
> machine but not on Gradescope this means you have a race condition that
> must be fixed in order for you to get full completeness points.

  - **Completeness:** 80%
  - **Correctness:** 10%
  - **Design & Style** 10%

## Obtaining your test score

The completeness part of your score will be determined using automated
tests. To get your score for the automated tests, simply run the
following from the **Terminal**. (Remember to leave out the `$` prompt
when you type the command.)

    $ cd grader
    $ go run hw2/grader hw2

This should print total score after running all test cases inside the
individual problems. This printout will not show the tests you failed.
You must run the problem's individual test file to see the failures.

### Design, Style and Cleaning up

Before you submit your final solution, you should, remove

  - any `Printf` statements that you added for debugging purposes and
  - all in-line comments of the form: "YOUR CODE HERE" and "TODO ..."
  - Think about your function decomposition. No code duplication. This
    homework assignment is relatively small so this shouldn't be a major
    problem but could be in certain problems.

Go does not have a strict style guide. However, use your best judgment
from prior programming experience about style. Did you use good variable
names? Do you have any lines that are too long, etc.

As you clean up, you should periodically save your file and run your
code through the tests to make sure that you have not broken it in the
process.

### Submission

Before submitting, make sure you've added, committed, and pushed all
your code to GitHub. You must submit your final work through Gradescope
(linked from our Canvas site) in the "Homework \#2" assignment page via
two ways,

1.  **Uploading from Github directly (recommended way)**: You can link
    your Github account to your Gradescope account and upload the
    correct repository based on the homework assignment. When you submit
    your homework, a pop window will appear. Click on "Github" and then
    "Connect to Github" to connect your Github account to Gradescope.
    Once you connect (you will only need to do this once), then you can
    select the repository you wish to upload and the branch (which
    should always be "main" or "master") for this course.
2.  **Uploading via a Zip file**: You can also upload a zip file of the
    homework directory. Please make sure you upload the entire directory
    and keep the initial structure the **same** as the starter code;
    otherwise, you run the risk of not passing the automated tests.

> **Note**
> For either option, you must upload the entire directory structure;
> otherwise, your automated test grade will not run correctly and you will
> be **penalized** if we have to manually run the tests. Going with the
> first option will do this automatically for you. You can always add
> additional directories and files (and even files/directories inside the
> stater directories) but the default directory/file structure must not
> change.

Depending on the assignment, once you submit your work, an "autograder"
will run. This autograder should produce the same test results as when
you run the code yourself; if it doesn’t, please let us know so we can
look into it. A few other notes:

  - You are allowed to make as many submissions as you want before the
    deadline.
  - Your completeness score is determined solely based on the automated
    tests, but we may adjust your score if you attempt to pass tests by
    rote (e.g., by writing code that hard-codes the expected output for
    each possible test input).
  - Gradescope will report the test score it obtains when running your
    code. If there is a discrepancy between the score you get when
    running our grader script, and the score reported by Gradescope,
    please let us know so we can take a look at it.
