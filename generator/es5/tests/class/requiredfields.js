$module('requiredfields', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      instance.SomeField = SomeField;
      instance.AnotherField = $t.box(true, $g.____testlib.basictypes.Boolean);
      return $promise.resolve(instance);
    };
    this.$typesig = function () {
      return $t.createtypesig(['new', 1, $g.____testlib.basictypes.Function($g.requiredfields.SomeClass).$typeref()]);
    };
  });

  $static.TEST = function () {
    var $result;
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.requiredfields.SomeClass.new($t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
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
            $g.____testlib.basictypes.Integer.$equals(sc.SomeField, $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result1) {
              return $promise.resolve($t.unbox($result1)).then(function ($result0) {
                $result = $t.box($result0 && $t.unbox(sc.AnotherField), $g.____testlib.basictypes.Boolean);
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
