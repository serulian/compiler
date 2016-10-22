$module('asyncstruct', function () {
  var $static = this;
  this.$struct('4e1a7940', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (Foo, Bar) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        Foo: Foo,
        Bar: Bar,
      };
      instance.$markruntimecreated();
      return $promise.resolve(instance);
    };
    $static.$fields = [];
    $t.defineStructField($static, 'Foo', 'Foo', function () {
      return $g.____testlib.basictypes.Integer;
    }, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'Bar', 'Bar', function () {
      return $g.____testlib.basictypes.Integer;
    }, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<4e1a7940>": true,
        "equals|4|29dc432d<5ab5941e>": true,
        "Stringify|2|29dc432d<538656f2>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<4e1a7940>": true,
        "String|2|29dc432d<538656f2>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomethingAsync = $t.workerwrap('dd7aa26e', function (s) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            s.Foo.String().then(function ($result1) {
              return s.Bar.String().then(function ($result2) {
                return $g.____testlib.basictypes.String.$plus($result1, $result2).then(function ($result0) {
                  $result = $result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                });
              });
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
  });
  $static.TEST = function () {
    var $result;
    var vle;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.asyncstruct.SomeStruct.new($t.fastbox(1, $g.____testlib.basictypes.Integer), $t.fastbox(2, $g.____testlib.basictypes.Integer)).then(function ($result1) {
              return $promise.translate($g.asyncstruct.DoSomethingAsync($result1)).then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            vle = $result;
            $g.____testlib.basictypes.String.$equals(vle, $t.fastbox("12", $g.____testlib.basictypes.String)).then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
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
