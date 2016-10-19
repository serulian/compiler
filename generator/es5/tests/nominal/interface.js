$module('interface', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
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
        $resolve($t.box(42, $g.____testlib.basictypes.Integer));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      var computed = $t.createtypesig(['SomeValue', 3, $g.____testlib.basictypes.Integer.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.interface.SomeClass).$typeref()]);
      this.$typesig = function () {
        return computed;
      };
      return computed;
    };
  });

  this.$interface('Valuable', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      var computed = $t.createtypesig(['SomeValue', 3, $g.____testlib.basictypes.Integer.$typeref()]);
      this.$typesig = function () {
        return computed;
      };
      return computed;
    };
  });

  this.$type('Valued', false, '', function () {
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
              $t.unbox($this).SomeValue().then(function ($result0) {
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
      var computed = $t.createtypesig(['GetValue', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Integer).$typeref()]);
      this.$typesig = function () {
        return computed;
      };
      return computed;
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
            $t.box(v, $g.interface.Valued).GetValue().then(function ($result1) {
              return $g.____testlib.basictypes.Integer.$equals($result1, $t.box(42, $g.____testlib.basictypes.Integer)).then(function ($result0) {
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
