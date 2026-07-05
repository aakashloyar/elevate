-> now the problem come here is ki for question and quiz 
-> do i need to create separate microservices or same microservice is fine for both
-> if they are belonging to same domain then same otherwise different
-> initially our plan is to just a quiz driven where questions are owned by assessment
-> but we can add functionality that same question belong to different assessment 
-> so no assessment does not own problem so different services


-> now the new problem is who own assessmentid problemid table 
-> accroding to DDD it must be assessment 
-> but problem is suppose user wants to access a quiz 
then he request question then from assessment service 
you will get which question then you request to problem service 
to give these problems is this a correct approach 
if it will be in problem service then we can directly 
get all problems there only

-> but in option2 if we think 
-> do problem service care about that talbe?
-> no -> so we need to compromise with a extra api call
-> for this we can use BFF


-> ok so let us suppose if i use this assessment owns 
that table then suppose for getting question of full assessment 
we will request to assessment service it will give list of all the problems ids 
now do i need to make different request for each problem 
or only 1 request and getting all the problems
-> we will try to avoid multiple request and use single request