1. POST /generation-jobs/{jobId}/cancel -> in future

2. not giving flexibility for option count
-> currently fixing to 4 will see in future if required

3. for topic we can just use json array or different table
we chooose different table 
it is no required as we donot search in. this service with relative to topic
but for future we have added it

# Gemini API Temperature
temperature parameter controls the degree of randomness or creativity in the model's token selection during text generation. 
How Temperature WorksLower values (closer to 0.0): Make the output more predictable, deterministic, and focused. 
The model consistently picks tokens with the highest probabilities.
Higher values (closer to 2.0): Increase randomness, leading to more creative, diverse, or unexpected responses.
Default value: Typically 1.0 for many models, which balances creativity and coherence.