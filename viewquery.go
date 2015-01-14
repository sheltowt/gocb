package gocouchbase

import "strconv"
import "encoding/json"

type StaleMode int

const (
	Before = StaleMode(1)
	None   = StaleMode(2)
	After  = StaleMode(3)
)

type SortOrder int

const (
	Ascending  = SortOrder(1)
	Descending = SortOrder(2)
)

type ViewQuery struct {
	ddoc    string
	name    string
	options map[string]string
}

func (vq *ViewQuery) Stale(stale StaleMode) *ViewQuery {
	if stale == Before {
		vq.options["stale"] = "false"
	} else if stale == None {
		vq.options["stale"] = "ok"
	} else if stale == After {
		vq.options["stale"] = "update_after"
	} else {
		panic("Unexpected stale option")
	}
	return vq
}

func (vq *ViewQuery) Skip(num uint) *ViewQuery {
	vq.options["skip"] = strconv.FormatUint(uint64(num), 10)
	return vq
}

func (vq *ViewQuery) Limit(num uint) *ViewQuery {
	vq.options["limit"] = strconv.FormatUint(uint64(num), 10)
	return vq
}

func (vq *ViewQuery) Order(order SortOrder) *ViewQuery {
	if order == Ascending {
		vq.options["descending"] = "false"
	} else if order == Descending {
		vq.options["descending"] = "true"
	} else {
		panic("Unexpected order option")
	}
	return vq
}

func (vq *ViewQuery) Reduce(reduce bool) *ViewQuery {
	if reduce == true {
		vq.options["reduce"] = "true"
	} else {
		vq.options["reduce"] = "false"
	}
	return vq
}

func (vq *ViewQuery) Group(level int) *ViewQuery {
	if level >= 0 {
		vq.options["group"] = "false"
		vq.options["group_level"] = strconv.FormatInt(int64(level), 10)
	} else {
		vq.options["group"] = "true"
		vq.options["group_level"] = "0"
	}
	return vq
}

func (vq *ViewQuery) Key(key string) *ViewQuery {
	jsonKey, _ := json.Marshal(key)
	vq.options["key"] = string(jsonKey)
	return vq
}

func (vq *ViewQuery) Keys(keys []string) *ViewQuery {
	jsonKeys, _ := json.Marshal(keys)
	vq.options["keys"] = string(jsonKeys)
	return vq
}

func (vq *ViewQuery) Range(start, end interface{}, inclusive_end bool) *ViewQuery {
	// TODO(brett19): Not currently handling errors due to no way to return the error
	if start != nil {
		jsonStartKey, _ := json.Marshal(start)
		vq.options["startkey"] = string(jsonStartKey)
	} else {
		delete(vq.options, "startkey")
	}
	if end != nil {
		jsonEndKey, _ := json.Marshal(end)
		vq.options["endkey"] = string(jsonEndKey)
	} else {
		delete(vq.options, "endkey")
	}
	if start != nil || end != nil {
		if inclusive_end {
			vq.options["inclusive_end"] = "true"
		} else {
			vq.options["inclusive_end"] = "false"
		}
	} else {
		delete(vq.options, "inclusive_end")
	}
	return vq
}

func (vq *ViewQuery) IdRange(start, end string) *ViewQuery {
	if start != "" {
		vq.options["startkey_docid"] = start
	} else {
		delete(vq.options, "startkey_docid")
	}
	if end != "" {
		vq.options["endkey_docid"] = end
	} else {
		delete(vq.options, "endkey_docid")
	}
	return vq
}

func (vq *ViewQuery) Custom(name, value string) *ViewQuery {
	vq.options[name] = value
	return vq
}

func NewViewQuery(ddoc, name string) *ViewQuery {
	return &ViewQuery{
		ddoc:    ddoc,
		name:    name,
		options: make(map[string]string),
	}
}
