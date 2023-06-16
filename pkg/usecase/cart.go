package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
	utils "github.com/rganes5/maanushi_earth_e-commerce/pkg/utils"
)

type CartUseCase struct {
	CartRepo interfaces.CartRepository
}

func NewCartUseCase(repo interfaces.CartRepository) services.CartUseCase {
	return &CartUseCase{
		CartRepo: repo,
	}
}

func (c *CartUseCase) AddToCart(ctx context.Context, productId string, id uint) error {
	//Find the corresponding cart with the user id and retreving the main cart.
	cart, err := c.CartRepo.FindCartById(ctx, id)
	// fmt.Println("cart is", cart)
	if err != nil {
		return err
	}

	product, err3 := c.CartRepo.FindProductById(ctx, productId)
	// fmt.Println("product is", product)
	if err3 != nil {
		return err3
	}

	//Finding the product details of the product which is being added to cart with passing the product id
	productDetails, err1 := c.CartRepo.FindProductDetailsById(ctx, productId)
	// fmt.Println("product details is", productDetails)
	if err1 != nil {
		return err1
	}

	if productDetails.InStock <= 0 {
		return errors.New("the item is out of stock")
	} else {
		fmt.Println("this is the productId from the usecase", productId)
		fmt.Println("this is the cartId from the usecase", cart.ID)
		//Finding whether the product is already added once so that we can update the quantity of the product and also update the total price of that specific product in cartItem table.
		existingItem, err2 := c.CartRepo.FindDuplicateProduct(ctx, productId, cart.ID)
		fmt.Println("Existing item returned from findduplicate function to usecase is", existingItem.ProductId)
		if err2 != nil {
			return err2

		}
		//Updating the existing cart details with quantity and price in database.
		if existingItem.ID != 0 {
			if productDetails.InStock <= existingItem.Quantity {
				return errors.New("the item is out of stock")
			}
			existingItem.Quantity++
			existingItem.TotalPrice = existingItem.Quantity * product.DiscountPrice
			err := c.CartRepo.UpdateCartItem(ctx, existingItem)
			// fmt.Println("updating cart item with", existingItem)
			if err != nil {
				return err
			}
		} else {
			if productDetails.InStock <= existingItem.Quantity {
				return errors.New("the item is out of stock")
			}
			//If the item being added is a new one, then we can populate a new struct with new item.
			newItem := domain.CartItem{
				CartID:     cart.ID,
				ProductId:  product.ID,
				Quantity:   1,
				TotalPrice: product.DiscountPrice,
			}
			//we can update the cart with the new item
			// fmt.Println("adding new item", newItem)
			if err := c.CartRepo.AddNewItem(ctx, newItem); err != nil {
				return err
			}

		}
	}
	return nil
}

func (c *CartUseCase) ListCart(ctx context.Context, id uint, pagination utils.Pagination) (int, []utils.ResponseCart, error) {
	//Find the corresponding cart with the user id and retreving the main cart.
	cart, err1 := c.CartRepo.FindCartById(ctx, id)
	cartItems, err2 := c.CartRepo.ListCart(ctx, id, pagination)
	err := errors.Join(err1, err2)
	if err != nil {
		return cart.GrandTotal, cartItems, err
	}
	return cart.GrandTotal, cartItems, err
}

func (c *CartUseCase) RemoveFromCart(ctx context.Context, productId string, id uint) error {
	//Find the corresponding cart with the user id and retreving the main cart.
	cart, err := c.CartRepo.FindCartById(ctx, id)
	// fmt.Println("cart is", cart)
	if err != nil {
		return err
	}

	product, err3 := c.CartRepo.FindProductById(ctx, productId)
	// fmt.Println("product is", product)
	if err3 != nil {
		return err3
	}
	existingItem, err2 := c.CartRepo.FindDuplicateProduct(ctx, productId, cart.ID)
	fmt.Println("Existing item returned from find duplicate function to usecase is", existingItem.ProductId)
	if err2 != nil {
		return err2
	}
	//Updating the existing cart details with quantity and price in database.
	if existingItem.ID != 0 {
		if existingItem.Quantity > 1 {
			existingItem.Quantity--
			existingItem.TotalPrice = existingItem.Quantity * product.DiscountPrice
			err := c.CartRepo.UpdateCartItem(ctx, existingItem)
			// fmt.Println("updating cart item with", existingItem)
			if err != nil {
				return err
			}
		} else if existingItem.Quantity == 1 {
			err4 := c.CartRepo.DeleteFromCart(ctx, existingItem)
			if err4 != nil {
				return err
			}
		}
	} else {
		return errors.New("the item does not exist in cart")
	}
	return nil
}
