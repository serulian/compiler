$module('interface', function () {
  var $static = this;
  this.$class('87868828', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return $promise.resolve(instance);
    };
    $instance.SomeValue = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.fastbox(42, $g.____testlib.basictypes.Integer));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeValue|3|c44e6c87": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$interface('0c7df867', 'Valuable', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "SomeValue|3|c44e6c87": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$type('76c61a88', 'Valued', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $g.interface.Valuable;
    };
    $instance.GetValue = function () {
      var $this = this;
      var $result;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        while (true) {
          switch ($current) {
            case 0:
              $this.$wrapped.SomeValue().then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
              }).catch(function (err) {
                $reject(err);
                return;
              });
              return;

            case 1:
              $resolve($result);
              return;

            default:
              $resolve();
              return;
          }
        }
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "GetValue|2|29dc432d<c44e6c87>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $result;
    var sc;
    var v;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.interface.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            sc = $result;
            v = sc;
            $t.fastbox(v, $g.interface.Valued).GetValue().then(function ($result1) {
              return $g.____testlib.basictypes.Integer.$equals($result1, $t.fastbox(42, $g.____testlib.basictypes.Integer)).then(function ($result0) {
                $result = $result0;
                $current = 2;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
