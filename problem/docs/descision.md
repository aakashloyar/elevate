1. first thing come in my mind is also like assesment do we have permission for question also 
-> but for now we must not think that much about future 
-> as now it is not the usecase so will see it in future
-> to make it futureproof we should not do overengineering 

2. picture allowed inside problem 
-> later 

3. problem metadata
-> later

4. make option reusable 
-> not good idea to save some data storge
-> we are searching for that option 
-> in a big data 



no i am just asking for 3 cases for choosing problem addition by ai
1. ki assessment will be handling this 
-> as problem are added for a specific assessment 
-> so assessment must handle this 
-> then it must send data llm 
-> which will then generate create quiz
-> then llm send events to problem and assessment service 

2. ki problem will be handling this
-> as problems will be added so it must be concern for problem 
-> then it must send data llm 
-> which will then generate create quiz
-> then llm send events to problem and assessment service 

3. ki llm only handle this 
-> which will then generate create quiz
-> then llm send events to problem and assessment service 
-> here the benefit is direct action no extra call 
-> as the call will not come by going to 1 extra service
-> but as this service is not owner or problem and assessment 
-> so they must only handle it

-> option 3 is preferred