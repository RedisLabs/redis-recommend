# redis-recommend
## A Simple Redis recommendation engine written in [Go](https://golang.org/).

###

## What is a recommendation engine?
A recommendation engine is a system that can predict what items each user would be interested in from given a set of items. 
Recommendation engines are used in a variety of applications ranging from e-commerce to online dating. 

## The algorithm
There are two basic approaches for building recommendation engines:
* **Content based classification:**

 This approach relies on classification by a large number of item and user attributes, assigning each user to possible classes of items.
  * pros: 

       - Can be very targeted 
	   - Allows detailed control to the system owner
	   - Does not require the user`s history.

  * cons:

       	- Requires deep knowledge of the items
	   	- Complicated data model
	   	- A lot of manual work to enter the items
		- Usually requires the user to enter a lot of details
	
	* would be best for: 
	
	   	- Dating
		- Resturant recommendations for lactose allergic people who like dairy food
	
* **Collaborative filtering:**

 This approach is based on taking user behaviour and making recommendations for the user based on users with similar behaviour. 
  * pros: 

		- Easy to implement with Redis!
		- Very generic, content of the item is irrelevant
		- Can yield unexpected relevant results
		

  * cons:
		
		- Requires a minimal level of a user's history before recommendation is viable
		- Can be computationaly heavy
		
	* Would be best for:
	
		- Movie/Music recommendations
		
As recommended by my local recommendation engine, I decided to go for a collaborative filtering example.

The algorithm is simple:

 * For a given user, find the top similar users by:
			
		1. Find all the users that rated at least one (or N) common item as the user and use them as candidates
		
		2. For each candidate, calculate a score using the RootMeanSquare of the difference between their mutual item ratings
		
		3. Store the top similar users for the user.
		
* Now find the top item recommendations by:
	
		1. Find all the items rated by the user`s top similars that the user *has not* rated

		2. Calculate the average rating for each item

		3. Store the top items
		

## The Redis

Each rating event of user U giving item I rating of R is stored in two sorted sets, the user's and the item's:

```
ZADD user:U:items R I
ZADD item:I:scores R U
```

Note that for our algorithm this is all the input data that we need!

To get user U items:

```
ZRANGE user:U:items 0 -1
```

In order to get the similarity candidates for user U we need the union of all the users that have mutually rated items with U. Let's assume U rated items I1, I2, I3:

```
ZUNIONSTORE ztmp 3 item:I1:scores item:I2:scores item:I3:scores
```

note: We stored the union in a temporary sorted set - "ztmp". 

Now let's use ZRANGE to fetch:

```
ZRANGE ztmp 0 -1
``` 

Now we need to calculate the similarity for each of the candidates. 
Assuming user U1 and U2, we want the RMS of all the diffs in the ratings of the item rated by both users.

Redis gives us ZINTERSTORE so we can get the intersection between U1 and U2 items. 

But how will we calculate the diff?  

This can be achieved by using weights.

Multiplying U1's ratings by -1 and U2's ratings by 1 will give us the diff:

```
ZINTERSTORE ztmp 2 user:U1:items user:U2:items WEIGHTS 1 -1  
```

After calculating the RMS in the client the results will be stored in user:U1:similars

Now that we have a sorted set of U1's similar users, we can extract the items that the similar users rated.

We'll do this with ZUNIONSTORE with all U1's similar users, but how can we make sure to exclude all the items U1 has already rated?

Weights are going to be used again, this time with the AGGREGATE option and ZRANGEBYSCORE command.

Multiplying U1's items by -1 and all the others by 1 and specifying the AGGREGATE MIN option will yield a sorted set that is easy to cut:

All U1's item scores will be negative while the other users item scores will be positive.

With ZRANGEBYSCORE we will fetch the items with a score > 0, giving us just what we wanted.

Assuming U1 with similar users U3,U5,U6:

```
ZUNIONSTORE ztmp 4 user:U1:items user:U3:items user:U5:items user:U6:items WEIGHTS -1 1 1 1 AGGREGATE MIN
ZRANGEBYSCORE ztmp 0 inf
```

The last step would be to calculate a score for each of the candidate items, which is the avarage rate given by U1's similar users.

To get all the rates of item I given user U1's similars we shall intersect the two sets and take only the item scores by using WEIGHTS:

```
ZINTERSTORE ztmp 2 user:U1:similars item:I:scores WEIGHTS 0 1 
```

## The Code

![alt text](http://blog.wiser.com/wp-content/uploads/2014/07/tumblr_lu7nekn38a1qfvq9bo1_500.jpg "Logo Title Text 1")


```go

	weights = append(weights, "WEIGHTS", -1.0)
	for _, simuser := range similarUsers {
		args = append(args, fmt.Sprintf("user:%s:items", simuser))
		weights = append(weights, 1.0)
	}
```	


