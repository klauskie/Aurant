# Aurant Web Services

# POST

## New Item
>/item

    {
        "rest_id":<number>,
        "name": <string>,
        "price": <string>
    }

## Item Update
>/item/update

    {
        "item_id": <number>,
        "name": <string>,
        "price": <string>
    }

## New Order
>/order

    {
        "item_id": <number>,
        "rest_id": <number>,
        "email": <string>
    }

## Update Order State
>/order/update

    {
        "cart_id": <number>,
        "state_id": <number>
    }

# GET

## Item Structure
>Get all Items -> /item

>Get Item by ID -> /item/detail/{number}

Returns list of objects in the following structure:

    [
        {
            "item_id": <number>,
            "rest_id": <number>,
            "category_id": <number>,
            "name": <string>,
            "description": <string>,
            "price": <string>,
            "is_enabled": <bool>
        }
    ]


## Order Structure
>Get all orders -> /order

>Get Orders by rest_id and state_id -> /order/restaurant/{number}/state/{number}

>Get Orders by Client -> /order/client/{email}/restaurant/{number}

Returns list of objects in the following structure:

    [
        {
            "cart_id": <number>,
            "rest_id": <number>,
            "item_id": <number>,
            "email": <string>,
            "datetime": "string|2019-08-24 23:39:10",
            "state_id": <number>,
            "additional_info": <string>,
            "full_item": {
                "item_id": <number>,
                "rest_id": <number>,
                "category_id": <number>,
                "name": "Tacos al Pastor",
                "description": <string>,
                "price": <string>,
                "is_enabled": <bool>
            }
        }
    ]


## Category Structure
>Get Categories by Restaurant ID -> /item/category/restaurant/{number}

Returns list of objects in the following structure:

    [
        {
            "category_id": <number>,
            "rest_id": <number>,
            "name": <string>
        }
    ]

## Restaurant Structure
>Get all restaurants -> /restaurant

Returns list of objects in the following structure:

    [
        {
            "rest_id": <number>,
            "name": <string>,
            "location": <string|lat,long>
        }
    ]

# Utilities 

## Increment Order State -> GET

>/update/increment/{number}

Returns:

    {
        "state": "2"
    }