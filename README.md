# Distributed-Systems-Synchro-DB


## A distributed database synchronizer made using Golang and RabbitMQ

<p align="center">
<img width="60%" src="https://github.com/yassinebk/Distributed-Systems-Synchro-DB/assets/65515933/38cb9a7c-5926-49f8-b26f-da7ce0b1a63c"/>
</p>


## The Problem

There are numerous methods available for synchronizing distributed databases. Consider the following scenario with specific constraints:

In the future system, there will be a main head office (HO) and multiple branch offices (BOs).
All offices are situated in different locations, and some of them face limitations in terms of internet connectivity. In certain cases, internet access may be available for only a few hours per day. To address this, we need to develop a custom database synchronization mechanism for data exchange between branches.


The BO sales branches are physically separated from the head office. They each manage their own databases independently, but it is necessary to synchronize their data with the head office, which maintains the complete sales data. We assume that the database is structured based on the product sales table with the following fields:

```
Product {
	ID          int
	Date        time.Time
	Product     string 
	Region      string
	Qty         uint32
	Cost        float32
	Tax         float32
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
```


## The Solution

Our database synchronization solution ensures a smooth data exchange between regional offices and the headquarters (HQ) through a structured process. To synchronize products, they are initially marked as "sent" and marshaled before being transmitted via a RabbitMQ queue. Upon receiving these messages, the HQ extracts and stores the products in its database, sending an acknowledgment (ack) for each stored product.

Subsequently, when the regional office receives the ack, it marks the product as both "sent" and "acked," avoiding duplicates and ensuring efficient synchronization between the regional offices and the HQ.

In cases where the ack is not received within a specified time frame, the regional office resets the product's status to "not sent," allowing for a reattempt at synchronization.

To handle potential edge cases, such as a RabbitMQ server crash before sending the ack, the HQ implements a validation step. Before syncing the product fields and returning the ack, the HQ checks if it has already stored a product with the corresponding external ID from the specific regional branch office. If the product exists, the HQ synchronizes the product fields and sends the ack, ensuring that the product is created only once in the HQ's database.




## Usage

![image](https://user-images.githubusercontent.com/62627838/236697266-eb0481c3-4fd5-4b4d-9f83-701fb813eb63.png)
![image](https://user-images.githubusercontent.com/62627838/236697296-92522112-9cfe-47be-ae50-6d31cc5ed6f5.png)
