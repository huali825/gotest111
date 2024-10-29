package repository

import (
	"context"
	"goworkwebook/webook003/internal/domain"
	"goworkwebook/webook003/internal/repository/cache"
	"goworkwebook/webook003/internal/repository/dao"
	"goworkwebook/webook003/pkg/logger"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error

	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)

	// BatchIncrReadCnt biz 和 bizId 长度必须一致 批量处理
	BatchIncrReadCnt(ctx context.Context, biz []string, bizId []int64) error
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
	l     logger.LoggerV1
}

func NewCachedInteractiveRepository(
	l logger.LoggerV1, dao dao.InteractiveDAO,
	cache cache.InteractiveCache) InteractiveRepository {
	return &CachedInteractiveRepository{dao: dao, cache: cache}
}

func (c *CachedInteractiveRepository) BatchIncrReadCnt(ctx context.Context, biz []string, bizId []int64) error {
	err := c.dao.BatchIncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	go func() {
		for i := 0; i < len(biz); i++ {
			er := c.cache.IncrReadCntIfPresent(ctx, biz[i], bizId[i])
			if er != nil {
				// 记录日志
			}
		}
	}()
	return nil
}

func (c *CachedInteractiveRepository) Get(
	ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	intr, err := c.cache.Get(ctx, biz, bizId)
	if err == nil {
		return intr, nil
	}

	daoIntr, err := c.dao.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	intr = c.toDomain(daoIntr)
	go func() {
		er := c.cache.Set(ctx, biz, bizId, intr)
		// 记录日志
		if er != nil {
			c.l.Error("回写缓存失败",
				logger.String("biz", biz),
				logger.Int64("bizId", bizId),
			)
		}
	}()
	return intr, nil
}

func (c *CachedInteractiveRepository) Liked(ctx context.Context,
	biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedInteractiveRepository) Collected(ctx context.Context,
	biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedInteractiveRepository) AddCollectionItem(ctx context.Context,
	biz string, id int64, cid int64, uid int64) error {
	err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Biz:   biz,
		BizId: id,
		Cid:   cid,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	return c.cache.IncrCollectCntIfPresent(ctx, biz, id)
}

func (c *CachedInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	// 你要更新缓存了
	// 部分失败问题 —— 数据不一致
	return c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

func (c *CachedInteractiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return c.cache.IncrLikeCntIfPresent(ctx, biz, id)
}

func (c *CachedInteractiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntIfPresent(ctx, biz, id)
}
func (c *CachedInteractiveRepository) toDomain(ie dao.Interactive) domain.Interactive {
	return domain.Interactive{
		ReadCnt:    ie.ReadCnt,
		LikeCnt:    ie.LikeCnt,
		CollectCnt: ie.CollectCnt,
	}
}
