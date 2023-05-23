package usecase

import (
	"context"
	"errors"

	"github.com/rganes5/maanushi_earth_e-commerce/pkg/domain"
	interfaces "github.com/rganes5/maanushi_earth_e-commerce/pkg/repository/interface"
	services "github.com/rganes5/maanushi_earth_e-commerce/pkg/usecase/interface"
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
	if err != nil {
		return err
	}
	//Finding the product and also the product details of the product which is being added to cart with passing the product id
	product, productDetails, err1 := c.CartRepo.FindProductDetailsById(ctx, productId)
	if err1 != nil {
		return err1
	}
	if productDetails.InStock <= 0 {
		return errors.New("the item is out of stock")
	} else {
		//Finding whether the product is already added once so that we can update the quantity of the product and also update the total price of that specific product in cartItem table.
		existingItem, err2 := c.CartRepo.FindDuplicateProduct(ctx, productId, cart.ID)
		if err2 != nil {
			return err2
		}
		//Updating the existing cart details with quantity and price in database.
		if existingItem.ID != 0 {
			existingItem.Quantity++
			existingItem.TotalPrice = existingItem.Quantity * product.Price
			if err := c.CartRepo.UpdateCartItem(ctx, existingItem); err != nil {
				return err
			} else {
				//If the item being added is a new one, then we can populate a new struct with new item.
				newItem := domain.CartItem{
					CartID:     cart.ID,
					ProductId:  product.ID,
					Quantity:   1,
					TotalPrice: product.DiscountPrice,
				}
				//we can update the cart with the new item
				if err := c.CartRepo.AddNewItem(ctx, newItem); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
