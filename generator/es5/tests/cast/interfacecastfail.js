$module('interfacecastfail', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.SomeValue = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box(2, $g.____testlib.basictypes.Integer));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      return $t.createtypesig(['SomeValue', 3, $g.____testlib.basictypes.Integer.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.interfacecastfail.SomeClass).$typeref()]);
    };
  });

  this.$interface('SomeInterface', false, '', function () {
    var $static = this;
    this.$typesig = function () {
      return $t.createtypesig(['SomeValue', 3, $g.____testlib.basictypes.Boolean.$typeref()]);
    };
  });

  $static.TEST = function () {
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.interfacecastfail.SomeClass.new().then(function ($result0) {
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
            $t.cast(sc, $g.interfacecastfail.SomeInterface);
            $resolve();
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