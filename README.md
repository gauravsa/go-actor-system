# Implementing Go Actor System

Actors are objects which encapsulate state and behavior, they communicate exclusively by exchanging messages which are placed into the recipient’s queue.
Read more about it [here](https://en.wikipedia.org/wiki/Actor_model).

Read more about this project implementation in [this article]https://medium.com/better-programming/implementing-the-actor-model-in-golang-3579c2227b5e

**Architecture**

![image info](./actor_system.png)


To use actor model:
Import the repo in your project and use it as done in [here](https://github.com/gauravsa/go-actor-system/blob/main/actor_system/actor_system_test.go)

<script src="https://gist.github.com/gauravsa/7a6feba10b06df259dab81697135aa96.js"></script>
